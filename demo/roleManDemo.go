package main
import "fmt"

func main() {
	fmt.Println(defaultRoleTypes())

	daddyRole := Role{roleName: "老爸", roleType: RoleTypeDaddy, roleGender: GenderMale, roleNick: "老爸", userNick: "哈哈👌"}

	fmt.Println(daddyRole)
	fmt.Println("角色名称:", daddyRole.roleName, "\n叫我：", daddyRole.userNick)

	changeRoleMommyType(&daddyRole)
	fmt.Println(daddyRole)
	changeRoleMommyTypeNoPointer(daddyRole)
	fmt.Println(daddyRole)


	waitingConfirmOrder := new(WaitingConfirmOrder)
	waitingConfirmOrder.orderId = "201903122323"
	waitingConfirmOrder.orderName = "淘宝商品"

	fmt.Println(waitingConfirmOrder.status())
	fmt.Println(waitingConfirmOrder.showDetail())
}



// 函数名

func defaultRoleTypes() []string {
	return []string{"老爸", "老妈", "男盆友", "女盆友", "儿子", "女儿"}
}


// 性别
const (
	GenderUnknown = iota
	GenderMale
	GenderFemale
)

// 角色类型
const (
	RoleTypeUnknown = iota
	RoleTypeDaddy
	RoleTypeMommy
	RoleTypeBoyFriend
	RoleTypeGirlFriend
	RoleTypeSon
	RoleTypeDaughter
)

const (
	OrderStatusWaitingConfirm = iota
	OrderStatusSettled
	OrderStatusInvalidate
	OrderStatusNoWard
)


// 结构体

type Role struct {
	roleName string
	roleType int
	roleGender int
	rolePoritrait string
	roleNick string
	userNick string
}



// 修改结构体的属性值

func changeRoleMommyType(role *Role) {
	role.roleName = "老妈"
	role.roleType = RoleTypeMommy
	role.roleNick = "老妈"
}


func changeRoleMommyTypeNoPointer(role Role) {
	role.roleName = "没有改变哦"
}




// 定义接口
type Order interface {
	status() int // 定义参数，类似getter & setter
	showDetail() string
}

// 待审核订单
type WaitingConfirmOrder struct {
	orderId string
	orderName string
}


// 实现接口
func (waitingConfirmOrder WaitingConfirmOrder) status() int {
	return OrderStatusWaitingConfirm
}

func (waitingConfirmOrder WaitingConfirmOrder) showDetail() string {
	result := "订单号：" + waitingConfirmOrder.orderId + "\n订单名称：" + waitingConfirmOrder.orderName

	return result
}

