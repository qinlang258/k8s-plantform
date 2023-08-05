package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type daemonSet struct{}

var DaemonSet daemonSet

type DaemonSetResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

func (d *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}

	return daemonSets
}

func (d *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (daemonSetResp *DaemonSetResp, err error) {
	daemonSetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: d.toCells(daemonSetList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Page:  page,
				Limit: limit,
			},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	daemonSets := d.fromCells(data.GenericDataList)
	return &DaemonSetResp{
		Items: daemonSets,
		Total: total,
	}, nil
}

// 获取daemonset详情
func (d *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error) {
	daemonSet, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return daemonSet, nil
}

func (d *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除daemonset失败，" + err.Error()))
		return errors.New("删除DaemonSet失败, " + err.Error())
	}
	return err
}

// 更新
func (d *daemonSet) UpdateDaemonSet(namespace, content string) error {
	var daemonSet = &appsv1.DaemonSet{}
	err := json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})

	if err != nil {
		logger.Error(errors.New("更新DaemonSet失败, " + err.Error()))
		return errors.New("更新DaemonSet失败, " + err.Error())
	}
	return nil
}
