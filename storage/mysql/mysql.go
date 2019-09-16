package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DeshErBojhaa/tradeshift/graph"
	// ...
	_ "github.com/go-sql-driver/mysql"
)

// MySQL ...
type MySQL struct {
	session *sql.DB
}

// CreateSchema bootstraps the initial database schema.
func (m *MySQL) CreateSchema() {
	log.Println("Creating schema fot tradeshift")
	m.session.Exec(`CREATE DATABASE IF NOT EXISTS tradeshift DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`)
	m.session.Exec(`USE tradeshift;`)
	m.session.Exec("CREATE TABLE IF NOT EXISTS nodes (Id varchar(20), ParId varchar(20) NULL, Height int)")
}

// NewMySQLStore creates an instance of MySQLStore with the given connection string.
func NewMySQLStore(connection string) (*MySQL, error) {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	// Check connection is up
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	session := &MySQL{session: db}
	session.CreateSchema()
	return session, nil
}

// InsertNode creates a graph node with it's parent child relationship into the datastore.
// Operations are done within a transaction to maintain data consistency.
func (m *MySQL) InsertNode(node *graph.Node) error {
	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	// Will not be called if commited prior
	defer tx.Rollback()

	stmtNode, err := tx.Prepare("INSERT INTO nodes (Id, ParId, Height) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtNode.Close()

	if _, err := stmtNode.Exec(node.ID, node.ParID, node.Height); err != nil {
		log.Printf("error happened executing node %#v", err)
		return err
	}

	return tx.Commit()
}

// GetNodes returns all the nodes. Intrensic info of a node is persisted in
// 'nodes' table. And the parent child relation is persisted in 'parents' table.
// First fetch all the info of the nodes. Then fetch all the info of parent child
// relation. Create child list for each node from that parent child relation.
// Return error at any point and avoid transaction.
func (m *MySQL) GetNodes() ([]*graph.Node, error) {
	nodes := make([]*graph.Node, 0)
	nodeMap := make(map[string]*graph.Node)

	rows, err := m.session.Query("SELECT Id, ParId, Height FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		node := graph.NewEmptyNode()
		rows.Scan(&node.ID, &node.ParID, &node.Height)
		nodes = append(nodes, &node)
		nodeMap[node.ID] = &node
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, node := range nodes {
		parNode := nodeMap[node.ParID]
		if parNode == nil { // Root
			continue
		}
		parNode.Children[node.ID] = node
	}

	return nodes, nil
}

// UpdateParent changes parent of 'curNode' to the 'targetNode'.
// If root is given as curNode an error is returned. It also updates
// the parent child relationship within a transaction.
func (m *MySQL) UpdateParent(curNode, targetNode *graph.Node) error {
	if curNode.ParID == "" {
		return fmt.Errorf("can not change parent of the root node")
	}
	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. All children of cur node should now be direct children of cur nodes parent (Move 1 level up)
	// 2. Cur node's parent will change

	// 1
	stmtLevelUpChildren, err := tx.Prepare("UPDATE nodes SET ParId=?, Height=Height-1 WHERE ParId=?")
	if err != nil {
		return err
	}
	defer stmtLevelUpChildren.Close()
	if _, err := stmtLevelUpChildren.Exec(curNode.ParID, curNode.ID); err != nil {
		return err
	}

	// 2
	stmtUpdatePar, err := tx.Prepare("UPDATE nodes SET ParId=?, Height=? WHERE Id=?")
	if err != nil {
		return err
	}
	defer stmtUpdatePar.Close()
	if _, err := stmtUpdatePar.Exec(targetNode.ID, targetNode.Height+1, curNode.ID); err != nil {
		return err
	}
	return tx.Commit()
}
