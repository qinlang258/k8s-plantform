package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

// dataselector 用于排序,过滤,分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell 接口,用于各种资源List的类型转换,转换后可以使用dataselector的排序,过滤,分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体,过滤:Name 分页:Limit和Page
type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

// FilterQuery 用于查询 过滤:Name
type FilterQuery struct {
	Name string
}

// 分页:Limit和Page Limit是单页的数据条数,Page是第几页
type PaginateQuery struct {
	Page  int
	Limit int
}

// Filter方法用于过滤,比较数据Name属性,若包含则返回
func (d *dataSelector) Filter() *dataSelector {
	if d.DataSelect.Filter.Name == "" {
		return d
	}
	filtered := []DataCell{}
	for _, value := range d.GenericDataList {
		// 定义是否匹配的标签变量,默认是匹配的
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue
		}
		if matches {
			filtered = append(filtered, value)
		}
	}
	d.GenericDataList = filtered
	return d
}

// Paginate 分页,根据Limit和Page的传参,取一定范围内的数据返回
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page
	//验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	//举例：25个元素的数组，limit是10，page是3，startIndex是20，endIndex是30（实际上endIndex是25）
	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 实现自定义的排序方法,需要重写Len,Swap,Less方法
// Len用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap用于数据比较大小后的位置变更
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less用于比较大小
func (d *dataSelector) Less(i, j int) bool {
	return d.GenericDataList[i].GetCreation().Before(d.GenericDataList[j].GetCreation())
}

// 重写以上三个方法,用sort.Sort 方法触发排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// deployCell
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func (d deploymentCell) GetName() string {
	return d.Name
}

type daemonSetCell appsv1.DaemonSet

func (d daemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d daemonSetCell) GetName() string {
	return d.Name
}

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.GetName()
}

type statefulSetCell appsv1.StatefulSet

func (s statefulSetCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s statefulSetCell) GetName() string {
	return s.Name
}

type nodeCell corev1.Node

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

type namespaceCell corev1.Namespace

func (n namespaceCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n namespaceCell) GetName() string {
	return n.Name
}
