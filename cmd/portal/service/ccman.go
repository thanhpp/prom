package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/thanhpp/prom/pkg/logger"

	"github.com/thanhpp/prom/pkg/etcdclient"
	"google.golang.org/grpc"

	"github.com/thanhpp/prom/pkg/timerpc"

	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/errconst"
)

// --------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CCMANService ----------------------------------------------------------

type CCManSrv struct {
	lock      sync.RWMutex
	clientMap map[string]ccmanrpc.CCManagerClient
	addr      map[string]string
}

func (cs *CCManSrv) client(shardID int) (client ccmanrpc.CCManagerClient, ok bool) {
	cs.lock.RLock()
	client, ok = cs.clientMap[fmt.Sprintf("cardscolumnsmanager-%d", shardID)]
	cs.lock.RUnlock()
	return client, ok
}

func (cs *CCManSrv) name() string {
	return "Card Column manager service"
}

func (cs *CCManSrv) respError(err error, respCode int32, respMessage string) error {
	return errconst.ServiceError{Srv: cs.name(), Err: err, Code: respCode, Msg: respMessage}
}

func (cs *CCManSrv) error(err error) error {
	return errconst.ServiceError{Srv: cs.name(), Err: err, Msg: err.Error()}
}

type iCCMan interface {
	ChooseShardIDFromInt(in int) (shardID int, err error)

	CreateCard(ctx context.Context, shardID int, card *ccmanrpc.Card) (err error)
	GetCardByID(ctx context.Context, shardID int, cardID uint32) (card *ccmanrpc.Card, err error)
	GetCardsByDueDate(ctx context.Context, shardID int, duedate time.Time) (cards []*ccmanrpc.Card, err error)
	GetCardsByAssignedToID(ctx context.Context, shardID int, userID uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByCreatorID(ctx context.Context, shardID int, userID uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByColumnID(ctx context.Context, shardID int, colID uint32) (cards []*ccmanrpc.Card, err error)
	UpdateCardByID(ctx context.Context, shardID int, cardID uint32, card *ccmanrpc.Card) (err error)
	DeleteCardByID(ctx context.Context, shardID int, cardID uint32) (err error)

	CreateColumn(ctx context.Context, shardID int, column *ccmanrpc.Column) (err error)
	GetColumnByID(ctx context.Context, shardID int, colID uint32) (col *ccmanrpc.Column, err error)
	GetColumnsByTitle(ctx context.Context, shardID int, title string) (cols []*ccmanrpc.Column, err error)
	GetColumnsByProjectID(ctx context.Context, shardID int, projectID uint32) (cols []*ccmanrpc.Column, err error)
	UpdateColumnByID(ctx context.Context, shardID int, colID uint32, col *ccmanrpc.Column) (err error)
	DeleteColumnByID(ctx context.Context, shardID int, colID uint32) (err error)
	DeleteColumnByIDAndMove(ctx context.Context, shardID int, colID uint32, newColID uint32) (err error)
}

var implCCManSrv = new(CCManSrv)

func SetCCManSrv(ctx context.Context) (err error) {
	implCCManSrv = &CCManSrv{
		clientMap: make(map[string]ccmanrpc.CCManagerClient),
		addr:      make(map[string]string),
	}

	services, err := etcdclient.Get().GetServices(ctx, "cardscolumnsmanager")
	if err != nil {
		return err
	}

	for i := range services {
		b := grpc.RoundRobin(etcdclient.Get().Resolver())
		conn, err := grpc.DialContext(ctx, services[i].Name, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancer(b))
		if err != nil {
			return err
		}
		implCCManSrv.lock.Lock()
		implCCManSrv.addr[services[i].Name] = services[i].Addr
		implCCManSrv.clientMap[services[i].Name] = ccmanrpc.NewCCManagerClient(conn)
		implCCManSrv.lock.Unlock()
	}

	return nil
}

func GetCCManSrv() iCCMan {
	return implCCManSrv
}

func checkServiceFailError(shardID int, err error) {
	if !strings.Contains(err.Error(), "code = Unavailable") {
		return
	}
	srv := fmt.Sprintf("cardscolumnsmanager-%d", shardID)
	addr, ok := implCCManSrv.addr[srv]
	if !ok {
		logger.Get().Error("Empty address")
	}
	// delete from map
	implCCManSrv.lock.Lock()
	delete(implCCManSrv.addr, srv)
	delete(implCCManSrv.clientMap, srv)
	implCCManSrv.lock.Unlock()
	// delete from etcd
	if err := etcdclient.Get().RemoveEndpoints(context.Background(), srv, addr); err != nil {
		logger.Get().Errorf("remove endpoints error: %v", err)
		checkServiceFailError(shardID, err)
		return
	}
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNCTIONS ----------------------------------------------------------

func (cS *CCManSrv) ChooseShardIDFromInt(in int) (shardID int, err error) {
	var (
		shardIDs []int
	)

	for k := range cS.clientMap {
		id, err := strconv.Atoi(strings.Split(k, "-")[1])
		if err != nil {
			return -1, err
		}
		shardIDs = append(shardIDs, id)
	}

	if len(shardIDs) == 0 {
		return -1, nil
	}

	return shardIDs[in%len(shardIDs)], nil

}

func (cS *CCManSrv) CreateCard(ctx context.Context, shardID int, card *ccmanrpc.Card) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.CreateCardReq{
		CreateCard: card,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.CreateCard(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) GetCardByID(ctx context.Context, shardID int, cardID uint32) (card *ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardByIDReq{
		CardID: cardID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetCardByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.ResCard, nil
}

func (cS *CCManSrv) GetCardsByDueDate(ctx context.Context, shardID int, duedate time.Time) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(err)
	}

	in := &ccmanrpc.GetCardsByDueDateReq{
		DueDate: timerpc.ToTimeRPC(duedate),
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetCardsByDueDate(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS *CCManSrv) GetCardsByAssignedToID(ctx context.Context, shardID int, userID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return
	}

	in := &ccmanrpc.GetCardsByAssignedToIDReq{
		AssignedToID: userID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetCardsByAssignedToID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS *CCManSrv) GetCardsByCreatorID(ctx context.Context, shardID int, userID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardsByCreatorIDReq{
		CreatorID: userID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetCardsByCreatorID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS *CCManSrv) GetCardsByColumnID(ctx context.Context, shardID int, colID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardsByColumnIDReq{
		ColumnID: colID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetCardsByColumnID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS *CCManSrv) UpdateCardByID(ctx context.Context, shardID int, cardID uint32, card *ccmanrpc.Card) (err error) {
	if ctx.Err() != nil {
		return cS.error(err)
	}

	in := &ccmanrpc.UpdateCardByIDReq{
		CardID:     cardID,
		UpdateCard: card,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.UpdateCardByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) DeleteCardByID(ctx context.Context, shardID int, cardID uint32) (err error) {
	if ctx.Err() != nil {
		return cS.error(cS.error(err))
	}

	in := &ccmanrpc.DeleteCardByIDReq{
		CardID: cardID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.DeleteCardByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) CreateColumn(ctx context.Context, shardID int, column *ccmanrpc.Column) (err error) {
	if ctx.Err() != nil {
		return
	}

	in := &ccmanrpc.CreateColumnReq{
		CreateColumn: column,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.CreateColumn(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) GetColumnByID(ctx context.Context, shardID int, colID uint32) (col *ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnByIDReq{
		ColumnID: colID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetColumnByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Column, nil
}

func (cS *CCManSrv) GetColumnsByTitle(ctx context.Context, shardID int, title string) (cols []*ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnsByTitleReq{
		Title: title,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetColumnsByTitle(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Columns, nil
}

func (cS *CCManSrv) GetColumnsByProjectID(ctx context.Context, shardID int, projectID uint32) (cols []*ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnsByProjectIDReq{
		ProjectID: projectID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return nil, fmt.Errorf("Client ID %d not found", shardID)
	}
	resp, err := client.GetColumnsByProjectID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return nil, cS.error(err)
	}

	return resp.Columns, nil
}

func (cS *CCManSrv) UpdateColumnByID(ctx context.Context, shardID int, colID uint32, col *ccmanrpc.Column) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.UpdateColumnByIDReq{
		ColumnID: colID,
		Column:   col,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.UpdateColumnByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) DeleteColumnByID(ctx context.Context, shardID int, colID uint32) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.DeleteColumnByIDReq{
		ColumnID: colID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("Client ID %d not found", shardID)
	}
	_, err = client.DeleteColumnByID(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}

func (cS *CCManSrv) DeleteColumnByIDAndMove(ctx context.Context, shardID int, colID uint32, newColID uint32) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.DeleteColumnByIDAndMoveReq{
		ColumnID:    colID,
		NewColumnID: newColID,
	}

	client, ok := cS.client(shardID)
	if !ok {
		return fmt.Errorf("DeleteColumnByIDAndMove %v", err)
	}

	_, err = client.DeleteColumnByIDAndMove(ctx, in)
	if err != nil {
		checkServiceFailError(shardID, err)
		return cS.error(err)
	}

	return nil
}
