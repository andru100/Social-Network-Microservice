import SendData from './SendData';

export default async function Chkauth (){ //check if user has valid jwt
    
  let jwt = localStorage.getItem('jwt_token');

  let jwtObject = {
    data: {
      Token: jwt
    }  
  }

  let gqlRequest = "query Chkauth ($data: JwtdataInput!){ Chkauth(input: $data) { AuthdUser } }"
  
  let response = await SendData(gqlRequest, jwtObject)

  if ( response ){ // if password is a match redirect to profile page
    console.log("authenticated username is:", response.data.Chkauth.AuthdUser)

    return response.data.Chkauth.AuthdUser
  } else {
    return false
  }
} 