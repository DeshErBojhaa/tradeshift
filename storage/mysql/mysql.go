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

	m.session.Exec("DROP TABLE nodes")
	m.session.Exec("DROP TABLE parents")
	m.session.Exec("CREATE TABLE nodes (Id varchar(20), ParId varchar(20), Height int)")
	m.session.Exec("CREATE TABLE parents (Id varchar(20), ParId varchar(20))")
}

// NewMySQLStore creates an instance of MySQLStore with the given connection string.
func NewMySQLStore(connection string) (*MySQL, error) {
	log.Println("Opening connection to:", connection)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	// Check connection is up
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// TODO: Create schema if does not exists
	return &MySQL{session: db}, nil
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

	stmtPar, err := tx.Prepare("INSERT INTO parents (Id, ParId) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmtPar.Close()

	if _, err := stmtNode.Exec(node.ID, node.ParID, node.Height); err != nil {
		return err
	}
	if _, err := stmtPar.Exec(node.ID, node.ParID); err != nil {
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
	childs := make(map[string][]string)

	rows, err := m.session.Query("SELECT Id, ParId, Height FROM nodes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		node := graph.Node{}
		rows.Scan(&node.ID, &node.ParID, &node.Height)
		nodes = append(nodes, &node)
		nodeMap[node.ID] = &node
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = m.session.Query("SELECT Id, ParId FROM parents")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var parentID, childID string
		rows.Scan(&childID, &parentID)
		childs[parentID] = append(childs[parentID], childID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, node := range nodes {
		childrenID := childs[node.ID]
		for _, cID := range childrenID {
			node.Children[cID] = nodeMap[cID]
		}
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

	stmtUpdatePar, err := tx.Prepare("UPDATE nodes SET ParId=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmtUpdatePar.Close()

	stmtUpdateRelation, err := tx.Prepare("UPDATE parents SET ParId=? WHERE id=? AND ParId=?")
	if err != nil {
		return err
	}
	defer stmtUpdateRelation.Close()

	if _, err := stmtUpdatePar.Exec(targetNode.ID, curNode.ID); err != nil {
		return err
	}

	if _, err := stmtUpdateRelation.Exec(targetNode.ID, curNode.ID, curNode.ParID); err != nil {
		return err
	}
	return tx.Commit()
}
