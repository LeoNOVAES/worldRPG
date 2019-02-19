package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Ids struct {
	ID int64 `json;"id"`
}

type Person struct {
	ID    string `json:"id,omitempty"`
	NAME  string `json:"name,omitempty"`
	RACA  string `json:"raca,omitempty"`
	CLASS string `json:"class,omitempty"`
	XP    string `json:"xp,omitempty"`
}

type Racas struct {
	ID   string `json:"id,omitempty"`
	RACA string `json:"raca,omitempty"`
}

type Classes struct {
	ID    string `json:"id,omitempty"`
	CLASS string `json:"class,omitempty"`
}

type Arms struct {
	ID   string `json:"id,omitempty"`
	ARMY string `json:"army,omitempty"`
}

type Power struct {
	ID    string `json:"id,omitempty"`
	POWER string `json:"power,omitempty"`
}

var db, err = sql.Open("mysql", "teste:123456@/world")

func auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Print("token verificado")
	} else {
		result := gin.H{
			"message": "nao autorizado",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}

}

func getArms(c *gin.Context) {
	id := c.Param("id")
	if err != nil {
		log.Fatal(err)
		c.JSON(500, err.Error())
	}

	rows, erro := db.Query("SELECT a.army FROM persons p JOIN ligation l ON l.id_person = p.id JOIN arms a ON l.id_arm = a.id WHERE p.id = ?", id)

	if erro != nil {
		log.Fatal(erro)
		c.JSON(500, erro.Error())
	}

	arms := make([]string, 0)

	for rows.Next() {
		var army string
		rows.Scan(&army)
		arms = append(arms, army)
	}

	c.JSON(200, arms)
	return

}

func cadastrarPerson(c *gin.Context) {
	Name := c.PostForm("name")
	Raca := c.PostForm("raca")
	Class := c.PostForm("class")
	Army := c.PostForm("arm")
	Power := c.PostForm("power")
	Xp := c.PostForm("xp")

	classF, err := strconv.ParseInt(Class, 10, 64)
	racaF, err := strconv.ParseInt(Raca, 10, 64)
	armyF, err := strconv.ParseInt(Army, 10, 64)
	powerF, err := strconv.ParseInt(Power, 10, 64)
	xpF, err := strconv.ParseInt(Xp, 10, 64)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, err.Error())
	}

	var nameLog string
	db.QueryRow("SELECT name FROM persons WHERE name = ?", Name).Scan(&nameLog)
	if len(nameLog) > 0 {
		c.JSON(400, "Usuario ja existe")
		return
	}

	stmt, er := db.Prepare("INSERT INTO persons(NAME,XP) VALUES(?,?)")

	if er != nil {
		log.Fatal(er)
		c.JSON(500, er.Error())
	}

	stmt.Exec(Name, xpF)

	id := 0
	db.QueryRow("SELECT id FROM persons WHERE name = ? ", Name).Scan(&id)

	cadLigation(id, racaF, classF, armyF, powerF)

	c.JSON(http.StatusOK, "Cadastrado com sucesso")
	return
}

func cadLigation(id int, raca int64, class int64, army int64, power int64) {

	if id > 0 {
		stmt, err := db.Prepare("INSERT INTO ligation(ID_PERSON, ID_RACA, ID_CLASS, ID_ARM, ID_POWER) VALUES(?,?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		stmt.Exec(id, raca, class, army, power)
	}
	return
}

func addArms(c *gin.Context) {
	idA := c.Param("id")
	army := c.PostForm("arm")

	id := getIdPersonForId(idA)

	if id == 0 {
		c.JSON(400, "Usuario nao existe")
		return
	}

	armyF, err := strconv.ParseInt(army, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	var numberArms int
	var idArm int64

	db.QueryRow("SELECT count(a.army), a.id  FROM persons p JOIN ligation l ON l.id_person = p.id JOIN arms a ON l.id_arm = a.id WHERE p.id = ?", id).Scan(&numberArms, &idArm)

	if numberArms >= 2 {
		c.JSON(400, "Nao pode cadastrar mais de duas armas")
		return
	}

	rows, err := db.Query("SELECT ligation.id_arm FROM ligation,persons  WHERE persons.id = ligation.id_person AND persons.id = ?", id)

	if err != nil {
		log.Fatal(err)
	}

	arm := 0

	for rows.Next() {
		i := Ids{}

		rows.Scan(&i.ID)

		idr := i.ID

		if idr == armyF {
			arm++
		}
	}

	if arm > 0 {
		c.JSON(400, "Arma ja existe para Personagem")
		return
	}

	stmt, _ := db.Prepare("INSERT INTO ligation(id_person, id_arm) VALUES(?,?)")
	stmt.Exec(id, armyF)

	c.JSON(200, "Arma cadastrada com sucesso")
	return

}

func getRacas(c *gin.Context) {
	rows, _ := db.Query("SELECT * FROM racas")

	var racas []Racas

	for rows.Next() {
		r := Racas{}
		rows.Scan(&r.ID, &r.RACA)
		racas = append(racas, r)
	}

	c.JSON(200, racas)
}

func getClasses(c *gin.Context) {
	rows, _ := db.Query("SELECT * FROM class")

	var class []Classes

	for rows.Next() {
		r := Classes{}
		rows.Scan(&r.ID, &r.CLASS)
		class = append(class, r)
	}

	c.JSON(200, class)
}

func getIdPersonForName(name string) int {
	var idF int
	db.QueryRow("SELECT id FROM persons WHERE name = ?", name).Scan(&idF)
	return idF
}

func getIdPersonForId(id string) int {
	var idF int
	db.QueryRow("SELECT id FROM persons WHERE id = ?", id).Scan(&idF)
	return idF
}

func addPower(c *gin.Context) {
	idP := c.Param("id")
	power := c.PostForm("power")
	id := getIdPersonForId(idP)

	if id == 0 {
		c.JSON(400, "Usuario nao existe")
		return
	}

	powerF, err := strconv.ParseInt(power, 10, 11)

	if err != nil {
		log.Fatal(err)
	}

	var name string
	var numbersId int
	var powerId int64
	var pv int = 0

	db.QueryRow("SELECT w.id,p.name, count(w.power) FROM persons p JOIN ligation l ON p.id = l.id_person JOIN powers w ON  w.id = l.id_power WHERE p.id = ?", id).Scan(&powerId, &name, &numbersId)

	if numbersId >= 5 {
		c.JSON(400, "Cada Personagem tem apenas 5 poderes")
		return
	}

	rows, errno := db.Query("SELECT ligation.id_power FROM ligation,persons  WHERE persons.id = ligation.id_person AND persons.id = ?", id)

	if errno != nil {
		log.Fatal(errno)

	}
	for rows.Next() {
		i := Ids{}
		rows.Scan(&i.ID)
		idd := i.ID

		if idd == powerF {
			pv++
		}
	}

	if pv > 0 {
		c.JSON(400, " "+name+" ja tem esse poder")
		return
	}

	stmt, _ := db.Prepare("INSERT INTO ligation(id_person, id_power) VALUES(?,?)")
	fmt.Print(id)
	stmt.Exec(id, powerF)
	c.JSON(200, "Poder cadastrado em "+name+" com sucesso")
	return
}

func getTotPowers(c *gin.Context) {
	rows, erro := db.Query("SELECT * FROM powers")

	if erro != nil {
		log.Fatal(erro)
	}

	var powers []Power

	for rows.Next() {
		p := Power{}
		er := rows.Scan(&p.ID, &p.POWER)
		if er != nil {
			log.Fatal(er)
		}
		powers = append(powers, p)
	}

	c.JSON(200, powers)
	return

}

func getTotArms(c *gin.Context) {
	rows, erro := db.Query("SELECT * FROM arms")

	if erro != nil {
		log.Fatal(erro)
	}
	var arms []Arms

	for rows.Next() {
		a := Arms{}
		er := rows.Scan(&a.ID, &a.ARMY)

		if er != nil {
			log.Fatal(er)
		}

		arms = append(arms, a)
	}
	c.JSON(200, arms)
	return
}

func getPersons(c *gin.Context) {

	if err != nil {
		log.Fatal(err)
		c.JSON(500, err.Error())
	}

	var persons []Person

	rows, erro := db.Query("SELECT p.id, p.name, c.raca, cc.class,p.xp FROM persons p JOIN ligation l ON l.id_person = p.id JOIN racas c ON l.id_raca = c.id  JOIN class cc ON cc.id = l.id_class GROUP BY p.name")

	if erro != nil {
		log.Fatal(erro)
		c.JSON(500, erro.Error())
	}

	for rows.Next() {
		p := Person{}
		er := rows.Scan(&p.ID, &p.NAME, &p.RACA, &p.CLASS, &p.XP)
		if er != nil {
			log.Fatal(er)
			c.JSON(500, er.Error())
		}
		persons = append(persons, p)
	}

	c.JSON(http.StatusOK, persons)
	return

}

func getPowerId(c *gin.Context) {
	id := c.Param("id")
	var powers []Power
	rows, _ := db.Query("SELECT w.power FROM persons p JOIN ligation l ON p.id = l.id_person JOIN powers w ON  w.id = l.id_power where p.id = ?", id)

	for rows.Next() {
		p := Power{}
		errno := rows.Scan(&p.POWER)
		if errno != nil {
			log.Fatal(err)
		}
		powers = append(powers, p)
	}

	c.JSON(200, powers)
	return

}

func deletePower(c *gin.Context) {
	idf := c.Param("id")
	power := c.Param("power")

	id, err := strconv.ParseInt(idf, 10, 8)
	powerf, err := strconv.ParseInt(power, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	w := 0
	rows, _ := db.Query("SELECT w.id FROM persons p JOIN ligation l ON p.id = l.id_person JOIN powers w ON  w.id = l.id_power where p.id = ?", id)
	var pw []Power
	for rows.Next() {

		p := Power{}
		er := rows.Scan(&p.ID)
		if er != nil {
			log.Fatal(er)
		}
		pw = append(pw, p)
		ids := p.ID
		i, _ := strconv.ParseInt(ids, 10, 11)

		if powerf == i {
			stm, _ := db.Prepare("UPDATE ligation SET id_power = null WHERE id_person = ? AND id_power = ?")
			stm.Exec(id, powerf)
			w++
		}
	}
	fmt.Print(w)
	if w != 0 {
		c.JSON(200, "Deletedo com sucesso")
		return
	} else {
		c.JSON(400, "Poder Nao existe para usuario")
		return
	}

}

func getPersonId(c *gin.Context) {
	id := c.Param("id")
	var p Person
	erro := db.QueryRow("SELECT p.id, p.name, c.raca, cc.class,p.xp FROM persons p JOIN ligation l ON l.id_person = p.id JOIN racas c ON l.id_raca = c.id  JOIN class cc ON cc.id = l.id_class WHERE p.id = ?", id).Scan(&p.ID, &p.NAME, &p.RACA, &p.CLASS, &p.XP)

	if erro != nil {
		fmt.Println(erro)
		c.JSON(400, "Usuario nao existe")
	}

	c.JSON(200, p)
	return
}

func editPerson(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	raca := c.PostForm("raca")
	class := c.PostForm("class")

	var nameVerify string
	db.QueryRow("SELECT name FROM persons WHERE name = ?", name).Scan(&nameVerify)

	if len(nameVerify) > 0 {
		c.JSON(400, " "+nameVerify+" ja esta sendo Usado")
		return
	}

	if len(name) > 0 {
		stmt, _ := db.Prepare("UPDATE persons SET name = ? WHERE id = ?")
		stmt.Exec(name, id)
	}

	if len(raca) > 0 {
		racaF, err := strconv.ParseInt(raca, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		stmtRaca, _ := db.Prepare("UPDATE ligation,persons SET ligation.id_raca = ? WHERE ligation.id_person = ?")
		stmtRaca.Exec(racaF, id)
	}

	if len(class) > 0 {
		classF, err := strconv.ParseInt(class, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		stmtClass, _ := db.Prepare("UPDATE ligation,persons SET ligation.id_class = ? WHERE ligation.id_person = ?")
		stmtClass.Exec(classF, id)
	}

	c.JSON(200, "Sucesso ao alterar")
	return
}

func uploadMultipart(c *gin.Context) {
	id := c.Param("id")
	file, header, err := c.Request.FormFile("file")
	data, _ := ioutil.ReadAll(file)

	t := time.Now()
	n := t.Format("20060102150405")
	timestampF := string(n)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./uploads/"+header.Filename+timestampF+".png", data, 0666)

	if err != nil {
		log.Fatal(err)
	}

	caminho := "/images/" + header.Filename + timestampF + ".png"

	var qtdImage int
	db.QueryRow("SELECT COUNT(path) FROM imgAvatar WHERE id_person = ?", id).Scan(&qtdImage)

	if qtdImage == 1 {
		stmtUpdate, _ := db.Prepare("UPDATE imgAvatar SET imgAvatar.path = ? WHERE imgAvatar.id_person = ?")
		stmtUpdate.Exec(caminho, id)
		c.JSON(200, "ALTERADO COM SUCESSO")
		return
	}

	stmt, _ := db.Prepare("INSERT INTO imgAvatar(path,id_person) VALUES(?,?)")
	stmt.Exec(caminho, id)

	c.JSON(200, "Avatar adicionado com sucesso")

}

func getImageAvatar(c *gin.Context) {
	id := c.Param("id")
	var path string
	db.QueryRow("SELECT  i.path from persons p JOIN imgAvatar i ON i.id_person = ? ", id).Scan(&path)
	c.Redirect(http.StatusMovedPermanently, path)

}

func deleteImgAvatar(c *gin.Context) {
	id := c.Param("id")
	stmt, _ := db.Prepare("DELETE FROM imgAvatar WHERE id_person = ?")
	stmt.Exec(id)
	c.JSON(202, "Imagem excluida com sucesso")
}

func isLogin(c *gin.Context) {
	reqEmail := c.PostForm("email")
	reqPass := c.PostForm("password")

	var email string
	var pass string

	db.QueryRow("SELECT email,password FROM USERS WHERE email = ? AND password = ?", reqEmail, reqPass).Scan(&email, &pass)

	if len(email) > 0 && len(pass) > 0 {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": email,
			"password": pass,
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {
			log.Fatal(err)

		}
		c.JSON(200, tokenString)

	} else {
		c.JSON(401, "Nao autorizado")

	}
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	idF := 0
	var name string
	db.QueryRow("SELECT id,name FROM persons WHERE id = ?", id).Scan(&idF, &name)

	if idF == 0 {

		c.JSON(400, "Usuario não existe")
		return
	}

	stmt, _ := db.Prepare("DELETE FROM persons WHERE persons.id = ?")
	stmt.Exec(id)

	stmt2, _ := db.Prepare("DELETE FROM ligation WHERE ligation.id_person = ?")
	stmt2.Exec(id)

	c.JSON(200, "usuario "+name+" Excluido com sucesso")
	return
}

func main() {
	router := gin.Default()

	router.Static("/images", "./uploads")

	router.GET("/persons", getPersons)
	router.GET("/person/:id", getPersonId)
	router.GET("/arms/:id", getArms)
	router.GET("/arms", getTotArms)
	router.GET("/powers", getTotPowers)
	router.GET("/powers/:id", getPowerId)
	router.GET("/racas", getRacas)
	router.GET("/classes", getClasses)
	router.GET("/avatar/:id", getImageAvatar)

	router.POST("/isLogin", isLogin)
	router.POST("/add/arm/:id", auth, addArms)
	router.POST("/cadastro", auth, cadastrarPerson)
	router.POST("/add/power/:id", auth, addPower)
	router.POST("/edit/person/:id", auth, editPerson)
	router.POST("/upload/:id", auth, uploadMultipart)
	router.POST("/teste/:id/:power", deletePower)

	router.DELETE("/delete/:id", auth, deletePerson)
	router.DELETE("/deleteAvatar/:id", auth, deleteImgAvatar)
	router.Run(":9000")

}
