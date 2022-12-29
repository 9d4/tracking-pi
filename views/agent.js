import superagent from "superagent";

const API_ROOT = `/api`;

const responseBody = (res) => ({
  res: res.body,
  raw: res,
});

const error = (err) => ({
  res: err.response.body,
  raw: err.response,
});

const requests = {
  del: (url) =>
      superagent.del(`${API_ROOT}${url}`).then(responseBody),
  get: (url) =>
      superagent.get(`${API_ROOT}${url}`).then(responseBody),
  put: (url, body) =>
      superagent
          .put(`${API_ROOT}${url}`, body)
          .then(responseBody),
  post: (url, body) =>
      superagent
          .post(`${API_ROOT}${url}`, body)
          .then(responseBody)
          .catch(error),
};

const Industries = {
  getAll: () => requests.get("/industries"),
  store: (body) => requests.post("/industries", body),
}

const Volunteers = {
  store: (body) => requests.post("/volunteers", body),
  storePhoto: (id, body) => requests.post(`/volunteers/${id}/photo`, body),
}

const Logs = {
  store: (body) => requests.post("/logs", body),
}

const agents = {
  Industries,
  Logs,
  Volunteers,
};

export default agents;
