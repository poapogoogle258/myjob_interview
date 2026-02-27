package ollama_test

import (
	"fmt"
	"testing"

	"github.com/poapogoogle258/myjob_interview/internel/clients/ollama"
)

type Job string

func (j Job) GetJobDescription() string {
	return string(j)
}
func TestGetSkillsRequestFromContent(t *testing.T) {

	job := Job(fmt.Sprintf(`
		title : Senior Backend Developer
		company : Collective Wisdom Co., Ltd. 

		------------------------------------------------------------------------

		UpPass เรากำลังขยายทีม! และเราต้องการ "พี่ใหญ่" สาย Backend ที่ไม่ได้แค่เขียน Code เก่ง แต่ต้องแม่นเรื่อง Infra และมีหัวใจเป็น Leader มาช่วยกันปั้นระบบของ UpPass ให้ Scale ไปอีกขั้น!

		🌟 สิ่งที่คุณจะได้ทำที่ UpPass:

		Technical Leadership: เป็นแกนหลักในการออกแบบและพัฒนา Backend โดยใช้ Django (Python) เน้นความ Clean, Secure และ Scalable

		Infrastructure Mastery: วางโครงสร้าง Cloud Infrastructure คุมระบบ Docker, Kubernetes และ CI/CD ให้เสถียร 24/7

		Team Growth: คุมทีม Developer, ทำ Code Review และวางมาตรฐานการทำงานที่ยอดเยี่ยมให้กับทีม

		Architectural Vision: ออกแบบระบบให้รองรับการเติบโตและแก้ปัญหาที่ซับซ้อนด้วย Best Practices



		🛠️ คุณสมบัติที่เรามองหา:

		✅ Python/Django Specialist: เข้าใจลึกถึงขั้นปรับแต่ง ORM และ Performance ได้

		✅ Infra Expert: แม่นเรื่อง Cloud, Containerization และระบบ Security

		✅ Leadership Experience: เคยคุมทีมหรือมีทักษะในการสื่อสารและบริหารจัดการคนได้ดี

		✅ Problem Solver: เจอ Bug ยากๆ หรือระบบล่มแล้วไม่ตกใจ รู้วิธีการแก้ปัญหาที่เป็นระบบ



		ทำไมต้องมาที่ UpPass?

		เงินเดือนสมน้ำสมเนื้อกับความสามารถ 

		ได้ทำงานใน Tech Stack ที่ทันสมัย และมีสิทธิ์ตัดสินใจในงานสถาปัตยกรรมระบบ วัฒนธรรมองค์กรแบบสมัยใหม่ เน้นผลลัพธ์ (Flexible Working)

		สังคม​Flat สุดๆ​ ที่นี่ไม่มีท่าน​ ไม่ทีคุณ​ (ยกเว้นตอนดึงGit เเล้วConflict)​คุยกันได้ทุกระดับเหมือนเพื่อนร่วมวงหมูกะทะ​ พร้อมรับฟังทุกไอเดียเทพของคุณ

		วันลาที่ไม่ได้จำกัด....ถ้างานเสร็จ​จะเสด็จถึงดาวอังคารก็ไม่ว่ากัน

		บุฟเฟ่ต์​ทุกเดือน​ เดือนละครั้ง​ กินให้ยับเเล้วกลับไปลุยCodeต่อ​ (หยอกๆ)

		ปล่อยของได้เต็มที่​ ไม่ต้องกลัวโดนดองไอเดีย​ ถ้าคุณมีSolution เทพๆหรืออยากลองท่าใหม่ๆDjangoหรือInfraจัดมาเลย! เราเน้นงานเสร็จเเละดี​ไม่เน้นพิธีรีตอง



		📍 สนใจร่วมทีมกับเรา:

		ส่ง Resume และ Profile เทพๆ ของคุณมาได้เลยที่: 📧 Email: hr@uppass.io 📝 โดยระบุหัวข้ออีเมล: "สมัครตำแหน่ง Senior Backend Developer - UpPass"

		มาสร้างสิ่งที่ยิ่งใหญ่ไปด้วยกันที่ UpPass นะคะ!`))

	result, err := ollama.GetSkillsRequestFromContent(job)
	if err != nil {
		t.Errorf(`%s`, err)
	}

	fmt.Println(result)

}
