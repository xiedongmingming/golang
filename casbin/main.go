package main

import (
	"fmt"
	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
	"github.com/casbin/casbin"
	"github.com/zupzup/casbin-http-role-example/authorization"
	"github.com/zupzup/casbin-http-role-example/model"
	"log"
	"net/http"
	"time"
)

// casbin
// 是一个用GO语言打造的轻量级开源访问控制框架
// https://github.com/casbin/casbin
// 目前在GITHUB开源
// CASBIN采用了元模型的设计思想支持多种经典的访问控制方案.如基于角色的访问控制RBAC、基于属性的访问控制ABAC等

// 主要特性包括:
// 支持自定义请求的格式.默认的请求格式为: {subject, object, action}
// 具有访问控制模型MODEL和策略POLICY两个核心概念
// 支持RBAC中的多层角色继承.不止主体可以有角色资源也可以具有角色
// 支持超级用户如ROOT或ADMINISTRATOR.超级用户可以不受授权策略的约束访问任意资源
// 支持多种内置的操作符如KEYMATCH方便对路径式的资源进行管理.如/FOO/BAR可以映射到/FOO*

// CASBIN不做的事情:
// 身份认证(即验证用户的用户名、密码)CASBIN只负责访问控制
// 应该有其他专门的组件负责身份认证然后由CASBIN进行访问控制二者是相互配合的关系
// 管理用户列表或角色列表.CASBIN认为由项目自身来管理用户、角色列表更为合适
// CASBIN假设所有策略和请求中出现的用户、角色、资源都是合法有效的



func loginHandler(users model.Users) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		user, err := users.FindByName(name)
		if err != nil {
			writeError(http.StatusBadRequest, "WRONG_CREDENTIALS", w, err)
			return
		}
		// setup session
		if err := session.RegenerateToken(r); err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		session.PutInt(r, "userID", user.ID)
		session.PutString(r, "role", user.Role)
		writeSuccess("SUCCESS", w)
	})
}

func logoutHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := session.Renew(r); err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess("SUCCESS", w)
	})
}
func currentMemberHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := session.GetInt(r, "userID")
		if err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess(fmt.Sprintf("User with ID: %d", uid), w)
	})
}
func memberRoleHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := session.GetString(r, "role")
		if err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess(fmt.Sprintf("User with Role: %s", role), w)
	})
}

func adminHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeSuccess("I'm an Admin!", w)
	})
}
func createUsers() model.Users {
	users := model.Users{}
	users = append(users, model.User{ID: 1, Name: "Admin", Role: "admin"})
	users = append(users, model.User{ID: 2, Name: "Sabine", Role: "member"})
	users = append(users, model.User{ID: 3, Name: "Sepp", Role: "member"})
	return users
}


func main() {

	//e := &casbin.Enforcer{}
	//
	//e.InitWithFile("examples/basic_model.conf", "examples/basic_policy.csv")
	//
	//sub := "alice"
	//obj := "data1"
	//act := "read"
	//
	//if e.Enforce(sub, obj, act) == true {
	//	fmt.Println("通过")
	//} else {
	//	fmt.Println("未通过")
	//}
	//
	//// 采用管理API进行权限的管理如获取一个用户所有的角色
	//
	//// roles := e.GetRoles("alice")

	// setup casbin auth rules

	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")

	if err != nil {
		log.Fatal(err)
	}

	// setup session store
	engine := memstore.New(30 * time.Minute)

	sessionManager := session.Manage(engine, session.IdleTimeout(30*time.Minute), session.Persist(true), session.Secure(true))

	// setup users
	users := createUsers()

	// setup routes
	mux := http.NewServeMux()

	mux.HandleFunc("/login", loginHandler(users))
	mux.HandleFunc("/logout", logoutHandler())
	mux.HandleFunc("/member/current", currentMemberHandler())
	mux.HandleFunc("/member/role", memberRoleHandler())
	mux.HandleFunc("/admin/stuff", adminHandler())

	log.Print("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", sessionManager(authorization.Authorizer(authEnforcer, users)(mux))))

}
