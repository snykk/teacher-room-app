package main_test

import (
	"a21hc3NpZ25tZW50/app/model"
	"a21hc3NpZ25tZW50/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	main "a21hc3NpZ25tZW50"
	repo "a21hc3NpZ25tZW50/app/repository"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Teacher Room - Deploy to Fly IO", func() {
	var teacherRepo repo.TeacherRepo
	os.Setenv("DATABASE_URL", "postgres://postgres:12345678@localhost:5432/my_db")
	db := config.NewDB()
	conn, err := db.Connect()
	Expect(err).ShouldNot(HaveOccurred())

	if err = conn.Migrator().DropTable("teachers"); err != nil {
		panic("failed droping table:" + err.Error())
	}

	BeforeEach(func() {
		teacherRepo = repo.NewTeacherRepo(conn)
		err := conn.AutoMigrate(&model.Teacher{})
		err = db.Reset(conn, "teachers")
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Repository", func() {
		Describe("Teacher repository", func() {
			When("add teachers data to teachers table database postgres", func() {
				It("should save data teacher to teachers table database postgres", func() {
					teacher := model.Teacher{
						Name:         "Aditira",
						FieldOfStudy: "Math",
						Age:          22,
					}
					err := teacherRepo.AddTeacher(teacher)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Teacher{}
					conn.Model(&model.Teacher{}).First(&result)
					Expect(result.Name).To(Equal(teacher.Name))
					Expect(result.FieldOfStudy).To(Equal(teacher.FieldOfStudy))
					Expect(result.Age).To(Equal(teacher.Age))

					err = db.Reset(conn, "teachers")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("read all data teacher from teachers table database postgres", func() {
				It("should return a list data teacher", func() {
					teacher := model.Teacher{
						Name:         "Aditira",
						FieldOfStudy: "Math",
						Age:          22,
					}
					err := teacherRepo.AddTeacher(teacher)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Teacher{}
					conn.Model(&model.Teacher{}).First(&result)
					Expect(result.Name).To(Equal(teacher.Name))
					Expect(result.FieldOfStudy).To(Equal(teacher.FieldOfStudy))
					Expect(result.Age).To(Equal(teacher.Age))

					res, err := teacherRepo.ReadTeacher()
					Expect(res[0].Name).To(Equal(teacher.Name))
					Expect(res[0].FieldOfStudy).To(Equal(teacher.FieldOfStudy))
					Expect(res[0].Age).To(Equal(teacher.Age))

					err = db.Reset(conn, "teachers")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update status teacher from teachers table database postgres", func() {
				It("should change field done true or false", func() {
					teacher := model.Teacher{
						Name:         "Aditira",
						FieldOfStudy: "Math",
						Age:          22,
					}
					err := teacherRepo.AddTeacher(teacher)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Teacher{}
					conn.Model(&model.Teacher{}).First(&result)
					Expect(result.Name).To(Equal(teacher.Name))
					Expect(result.FieldOfStudy).To(Equal(teacher.FieldOfStudy))
					Expect(result.Age).To(Equal(teacher.Age))

					err = teacherRepo.UpdateName(1, "Dito")
					Expect(err).ShouldNot(HaveOccurred())

					result = model.Teacher{}
					conn.Model(&model.Teacher{}).First(&result)
					Expect(result.Name).To(Equal("Dito"))

					err = db.Reset(conn, "teachers")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete data teacher from teachers table database postgres", func() {
				It("should remove data teacher from teachers table database postgres according to target", func() {
					teacher := model.Teacher{
						Name:         "Aditira",
						FieldOfStudy: "Math",
						Age:          22,
					}
					err := teacherRepo.AddTeacher(teacher)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Teacher{}
					conn.Model(&model.Teacher{}).First(&result)
					Expect(result.Name).To(Equal(teacher.Name))
					Expect(result.FieldOfStudy).To(Equal(teacher.FieldOfStudy))
					Expect(result.Age).To(Equal(teacher.Age))

					err = teacherRepo.DeleteTeacher(1)
					Expect(err).ShouldNot(HaveOccurred())

					resTeacher := model.Teacher{}
					conn.Table("teachers").Select("*").Scan(&resTeacher)
					Expect(resTeacher.DeletedAt.Valid).To(Equal(true))

					err = db.Reset(conn, "teachers")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("Deploy to fly.io", func() {
		Describe("Teacher", func() {
			When("base url fly.io", func() {
				It("no error returns", func() {
					req, err := http.NewRequest("GET", main.FlyURL(), nil)
					Expect(err).ShouldNot(HaveOccurred())

					var client = &http.Client{}
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())

					_, err = ioutil.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("/api/teacher/add", func() {
				It("no error returns", func() {
					postBody, _ := json.Marshal(map[string]interface{}{
						"name":           "Aditira",
						"field_of_study": "Math",
						"age":            22,
					})
					requestBody := bytes.NewBuffer(postBody)

					req, err := http.NewRequest("POST", main.FlyURL()+"/api/teacher/add", requestBody)
					Expect(err).ShouldNot(HaveOccurred())

					var client = &http.Client{}
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())

					var success model.SuccessResponse
					err = json.Unmarshal(body, &success)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(success.Message).To(Equal("Teacher Added"))

					_, err = http.Get(main.FlyURL() + "/api/teacher/reset")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("/api/teacher/read", func() {
				It("no error returns", func() {
					postBody, _ := json.Marshal(map[string]interface{}{
						"name":           "Aditira",
						"field_of_study": "Math",
						"age":            22,
					})
					requestBody := bytes.NewBuffer(postBody)

					_, err := http.Post(main.FlyURL()+"/api/teacher/add", "application/json", requestBody)
					Expect(err).ShouldNot(HaveOccurred())

					req, err := http.NewRequest("GET", main.FlyURL()+"/api/teacher/read", nil)
					Expect(err).ShouldNot(HaveOccurred())

					var client = &http.Client{}
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())

					var teachers []model.ViewTeacher
					err = json.Unmarshal(body, &teachers)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(teachers[0].Name).To(Equal("Aditira"))
					Expect(teachers[0].FieldOfStudy).To(Equal("Math"))
					Expect(teachers[0].Age).To(Equal(22))

					_, err = http.Get(main.FlyURL() + "/api/teacher/reset")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("/api/teacher/update", func() {
				It("no error returns", func() {
					postBody, _ := json.Marshal(map[string]interface{}{
						"name":           "Aditira",
						"field_of_study": "Math",
						"age":            22,
					})
					requestBody := bytes.NewBuffer(postBody)

					_, err := http.Post(main.FlyURL()+"/api/teacher/add", "application/json", requestBody)
					Expect(err).ShouldNot(HaveOccurred())

					postBody, _ = json.Marshal(map[string]interface{}{
						"id":       1,
						"new_name": "Dion",
					})
					requestBody = bytes.NewBuffer(postBody)

					req, err := http.NewRequest("PUT", main.FlyURL()+"/api/teacher/update", requestBody)
					Expect(err).ShouldNot(HaveOccurred())

					var client = &http.Client{}
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())

					var success model.SuccessResponse
					err = json.Unmarshal(body, &success)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(success.Message).To(Equal("Teacher Name Changed"))

					_, err = http.Get(main.FlyURL() + "/api/teacher/reset")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("/api/teacher/delete?id=1", func() {
				It("no error returns", func() {
					postBody, _ := json.Marshal(map[string]interface{}{
						"name":           "Aditira",
						"field_of_study": "Math",
						"age":            22,
					})
					requestBody := bytes.NewBuffer(postBody)

					_, err := http.Post(main.FlyURL()+"/api/teacher/add", "application/json", requestBody)
					Expect(err).ShouldNot(HaveOccurred())

					req, err := http.NewRequest("DELETE", main.FlyURL()+"/api/teacher/delete?id=1", nil)
					Expect(err).ShouldNot(HaveOccurred())

					var client = &http.Client{}
					resp, err := client.Do(req)
					Expect(err).ShouldNot(HaveOccurred())

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())

					var success model.SuccessResponse
					err = json.Unmarshal(body, &success)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(success.Message).To(Equal("Teacher Delete Success"))

					_, err = http.Get(main.FlyURL() + "/api/teacher/reset")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})
})
