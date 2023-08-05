package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type node struct {
}

var Node node

type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

func (n *node) toCells(std []corev1.Node) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range cells {
		cells[i] = nodeCell(std[i])
	}
	return cells
}

func (n *node) fromCells(cells []DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		nodes[i] = corev1.Node(cells[i].(nodeCell))
	}
	return nodes
}

func (n *node) GetNodes(filterName string, page, limit int) (nodeResp *NodeResp, err error) {
	nodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取容器失败" + err.Error()))
		return nil, errors.New("获取容器失败" + err.Error())
	}
	dataSelect := dataSelector{
		GenericDataList: n.toCells(nodeList.Items),
		DataSelect: &DataSelectQuery{
			Filter:   &FilterQuery{filterName},
			Paginate: &PaginateQuery{Page: page, Limit: limit},
		},
	}

	filtered := dataSelect.Filter()
	total := len(dataSelect.GenericDataList)
	data := filtered.Sort().Paginate()

	nodes := n.fromCells(data.GenericDataList)
	return &NodeResp{
		Items: nodes,
		Total: total,
	}, nil
}

func (n *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	item, err := K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return item, nil
}
