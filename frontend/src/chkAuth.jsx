export default async function Chkauth (){ //check if user has valid jwt
    let jwt = localStorage.getItem('jwt_token');

    let jwtObject = {
      Data1: jwt
    }

    let options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(jwtObject),
      }
      
      let url = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_CHKAUTH_PORT + '/chkauth'
      console.log('env var urls is', url)
      let response = await fetch(url, options)
      let convert = await response.json ()

    if ( response.status === 401 || response.status === 400){
      return false   
     } else if ( response.status === 200){ // if password is a match return auth'd username
      return convert.AuthdUser
     }
} 