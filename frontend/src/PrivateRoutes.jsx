import React from 'react';
import { Routes, Route} from 'react-router-dom';
import ProfileSetup from "./routes/ProfileSetup";


//Allows access to pages that require authorisation and passes the auth'd username in props
const PrivateRoute = (props) => {
  alert("private route called")
  return (
    <Routes>
      
      {/* <Route path="/" element={<ProfileSetup user={props.user}/>} /> */}
      <Route path="/profileSetup" element={<ProfileSetup user={props.user}/>} />
    </Routes>
  )
}

export default PrivateRoute;