package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceResp struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

var Namespace namespace

type namespace struct {
}

func (n *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range cells {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}

func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i := range namespaces {
		namespaces[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return namespaces
}

func (n *namespace) GetNamespaceDetail(namespaceName string) (namespace *corev1.Namespace, err error) {
	item, err := K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return item, nil
}

func (n *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func (n *namespace) GetNamespaces(filterName string, limit, page int) (namespaceResp *NamespaceResp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	dataSelect := &dataSelector{
		GenericDataList: n.toCells(namespaceList.Items),
		DataSelect:      &DataSelectQuery{Paginate: &PaginateQuery{Limit: limit, Page: page}, Filter: &FilterQuery{filterName}},
	}

	filtered := dataSelect.Filter()
	total := len(dataSelect.GenericDataList)
	data := filtered.Sort().Paginate()

	namespaces := n.fromCells(data.GenericDataList)
	return &NamespaceResp{Items: namespaces, Total: total}, nil
}
