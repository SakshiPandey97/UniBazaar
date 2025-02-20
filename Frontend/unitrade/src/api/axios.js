import axios from 'axios';

const USER_BASE_URL="http://127.0.0.1:4000"

export const userLoginAPI = ({userLoginObject}) => {
  return axios.post(USER_BASE_URL+'/login', 
    userLoginObject)
  .then(response => {
    return response.data.userId; 
  })
  .catch(error => {
    console.error('Error logging in:', error);
    throw error;
  });
};

export const userRegisterAPI = ({userRegisterObject}) => {
  console.log("Register Recevied Obj",userRegisterObject)
  return axios.post(USER_BASE_URL+'/signup', 
    userRegisterObject)
  .then(response => {
    return response.data; 
  })
  .catch(error => {
    console.error('Error Registering: ', error);
    throw error;
  });
};
