import "./base.scss";
import "bootstrap";
import agents from "./agent.js";

const cardGeo = document.getElementById("card_geo");
const inputCode = document.getElementById("input_code");
const cardSubmit = document.getElementById("card_submit");
const btnSubmit = document.getElementById("btn_submit");
const btnGetLoc = document.getElementById("btn_get_loc");
const camPlayer = document.getElementById('camera_player');
const camPwrBtn = document.getElementById("camera_power");
const canvas = document.getElementById('canvas');
const canvasContext = canvas.getContext('2d');
const btnFlipCam = document.getElementById("btn_flip");


let location = {
  latitude: null,
  longitude: null,
}
let photoUri = null;
let gotLocation = false;

function onGetlocSuccess(position) {
  location.latitude = position.coords.latitude;
  location.longitude = position.coords.longitude;
  gotLocation = true;
  cardGeo.innerHTML =
    `<div class="text-success">Arigatōgozaimashita</div>`

  cardGeo.classList.add("opacity-0");

  setTimeout(() => {
    cardGeo.classList.add("invisible");
  }, 1400);
}

function check() {
  if (inputCode.value !== "" && location.longitude !== null & location.latitude !== null && photoUri != null && gotLocation) {
    cardSubmit.classList.remove("d-none");
    return
  }

  cardSubmit.classList.add("d-none");
}

function onGetlocError() {
  cardGeo.innerHTML =
    `<div class="alert alert-danger mb-0">Waduuh... Kamu tidak memberikan izin padaku. Sepertinya kamu harus meresetku.</div>`
  cardGeo.classList.remove(["invisible", "opacity-0"]);
}

function whereAmI() {
  if (!navigator.geolocation) {
    onGetlocError();
  } else {
    navigator.geolocation.getCurrentPosition(onGetlocSuccess, onGetlocError);
  }
}

camPwrBtn.addEventListener("click", () => {
  const constraints = {
    video: true,
  };

  navigator.mediaDevices.getUserMedia(constraints).then((stream) => {
    // Attach the video stream to the video element and autoplay.
    camPlayer.srcObject = stream;
    camPlayer.classList.remove("d-none");
    camPwrBtn.classList.add("d-none");
  });

})

document.addEventListener("DOMContentLoaded", () => {
  setInterval(check, 1500);
});

const btnShoot = document.getElementById("btn_shoot");
const btnReset = document.getElementById("btn_reset");
btnShoot.addEventListener("click", () => {
  canvas.height = camPlayer.offsetHeight;
  canvas.width = camPlayer.offsetWidth;

  canvasContext.drawImage(camPlayer, 0, 0, canvas.width, canvas.height);
  photoUri = canvas.toDataURL("image/jpeg", 100);
  camPlayer.classList.add("d-none");
  canvas.classList.remove("d-none");
});

btnReset.addEventListener("click", () => {
  camPlayer.classList.remove("d-none");
  canvas.classList.add("d-none");
  photoUri = null;
})

btnFlipCam.addEventListener("click", () => {
  if (camPlayer.style.transform === "scaleX(-1)") {
    camPlayer.style.transform = "scaleX(1)"
    canvas.style.transform = "scaleX(1)"
  } else {
    camPlayer.style.transform = "scaleX(-1)"
    canvas.style.transform = "scaleX(-1)"
  }
})

btnSubmit.addEventListener("click", () => {
  let volunteer_code = inputCode.value;
  let coordinate = location
  let photo = photoUri;


  agents.Logs.store({
    volunteer_code,
    coordinate,
    photo,
  }).then(({res, raw}) => {
    if (raw.status !== 201) {
      alert("Terjadi kesalahan, Coba Lagi ya! Error:" + raw.statusText);
      return;
    }

    alert("Terimakasih! Datamu sudah tersimpan, kamu boleh meninggalkan laman ini.");
  })
})

btnGetLoc.addEventListener("click", () => {
  whereAmI();
})
