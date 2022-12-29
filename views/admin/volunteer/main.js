import "../../base.scss";
import "bootstrap";
import agents from "../../agent.js";

const newVolunteerEl = document.getElementById("link_new");
const newFormCardEl = document.getElementById("form_new_card");
const formNewEl = document.getElementById("form_new");
const formSubmitEl = document.getElementById("btn_submit");
const photoSelectEl = document.getElementById("input_photo");

newVolunteerEl.addEventListener("click", () => {
  newFormCardEl.classList.remove("visually-hidden");
})

formNewEl.addEventListener("submit", (ev) => ev.preventDefault());


const photoReader = new FileReader();
const photoFile = {
  dom: photoSelectEl,
  binary: null,
};

photoReader.addEventListener("load", () => {
  photoFile.binary = photoReader.result;
});

photoSelectEl.addEventListener("change", () => {
  if (photoReader.readyState === FileReader.LOADING) {
    photoReader.abort();
  }

  photoReader.readAsBinaryString(photoFile.dom.files[0]);
});

formSubmitEl.addEventListener("click", () => {
  let name = document.getElementById("input_name").value;
  let industry_code = document.getElementById("input_code").value;
  let photo = document.getElementById("input_photo").value;

  agents.Volunteers.store({
    name,
    industry_code
  }).then(({res, raw}) => {
    if (res === null) {
      alert(raw.statusText);
      return;
    }

    if (raw.status === 201) {
      // now upload photo
      // const formData = new FormData()
      // formData.append("photo", photo);

      sendPhoto(`/api/volunteers/${res.id}/photo`);
    }
  })
})

function sendPhoto(url) {
  // If there is a selected file, wait until it is read
  // If there is not, delay the execution of the function
  if (!photoFile.binary && photoFile.dom.files.length > 0) {
    setTimeout(sendPhoto, 10);
    return;
  }

  // To construct our multipart form data request,
  // We need an XMLHttpRequest instance
  const XHR = new XMLHttpRequest();

  // We need a separator to define each part of the request
  const boundary = "blob";

  // Store our body request in a string.
  let data = "";

  // So, if the user has selected a file
  if (photoFile.dom.files[0]) {
    // Start a new part in our body's request
    data += `--${boundary}\r\n`;

    // Describe it as form data
    data += 'content-disposition: form-data; '
      // Define the name of the form data
      + `name="${photoFile.dom.name}"; `
      // Provide the real name of the file
      + `filename="${photoFile.dom.files[0].name}"\r\n`;
    // And the MIME type of the file
    data += `Content-Type: ${photoFile.dom.files[0].type}\r\n`;

    // There's a blank line between the metadata and the data
    data += '\r\n';

    // Append the binary data to our body's request
    data += photoFile.binary + '\r\n';
  }

  // // Text data is simpler
  // // Start a new part in our body's request
  // data += `--${boundary}\r\n`;
  //
  // // Say it's form data, and name it
  // data += `content-disposition: form-data; name="${photoFile.dom.name}"\r\n`;
  // // There's a blank line between the metadata and the data
  // data += '\r\n';
  //
  // // Append the text data to our body's request
  // data += text.value + "\r\n";

  // Once we are done, "close" the body's request
  data += `--${boundary}--`;

  // Define what happens on successful data submission
  XHR.addEventListener('load', () => {
    window.location.reload();
  });

  // Define what happens in case of an error
  XHR.addEventListener('error', () => {
    alert('Oops! Something went wrong.');
  });

  // Set up our request
  XHR.open('POST', url);

  // Add the required HTTP header to handle a multipart form data POST request
  XHR.setRequestHeader('Content-Type', `multipart/form-data; boundary=${boundary}`);

  // Send the data
  XHR.send(data);
}
