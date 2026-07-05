package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/sauravsinghs/simplebank/db/sqlc"
	"github.com/sauravsinghs/simplebank/pb"
	"github.com/sauravsinghs/simplebank/token"
	"github.com/sauravsinghs/simplebank/util"
	"github.com/sauravsinghs/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}