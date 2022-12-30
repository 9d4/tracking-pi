import "../../base.scss";
import "bootstrap";
import agents from "../../agent.js";

const newVolunteerEl = document.getElementById("link_new");
const newFormCardEl = document.getElementById("form_new_card");
const formNewEl = document.getElementById("form_new");
const formSubmitEl = document.getElementById("btn_submit");
const photoSelectEl = document.getElementById("input_photo");

const convertBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const fileReader = new FileReader();
    fileReader.readAsDataURL(file);

    fileReader.onload = () => {
      resolve(fileReader.result);
    };

    fileReader.onerror = (error) => {
      reject(error);
    };
  });
};

newVolunteerEl.addEventListener("click", () => {
  newFormCardEl.classList.remove("visually-hidden");
})

formNewEl.addEventListener("submit", (ev) => ev.preventDefault());

let photoB64 = "";

const uploadPhoto = (id) => {
  if (photoB64 === "") {
    setTimeout(uploadPhoto, 1000);
    return;
  }

  agents.Volunteers.storePhoto(id, {
    photo: photoB64
  }).then(({res, raw}) => {
    if (raw.status === 201) {
      window.location.reload();
      return;
    }

    alert(raw.statusText);
  })
}

photoSelectEl.addEventListener("change", async (event) => {
  const file = event.target.files[0];
  photoB64 = await convertBase64(file);
});

formSubmitEl.addEventListener("click", () => {
  let name = document.getElementById("input_name").value;
  let industry_code = document.getElementById("input_code").value;
  // code is volunteer code
  let code = document.getElementById("input_code_volunteer").value;

  agents.Volunteers.store({
    name,
    code,
    industry_code
  }).then(({res, raw}) => {
    if (res === null) {
      alert(raw.statusText);
      return;
    }

    if (raw.status === 201) {
      // now upload photo
      uploadPhoto(res.id);
    }
  })
})
