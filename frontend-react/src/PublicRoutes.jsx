import React from 'react';
import { Routes, Route} from 'react-router-dom';
import SignUp from "./routes/SignUp";
import SignIn from "./routes/SignIn";
import ProfileSetup from "./routes/ProfileSetup";
import Profile from "./routes/Profile";
import UpdateDetails from "./routes/UpdateDetails";
import ResetComplete from "./routes/ResetComplete";
import Verify from "./routes/Verify";

//Allows access to pages that have no authorisation and allows them to sign in/up

const PublicRoute = () => {
  return (
    <Routes>
      <Route path="/" element={<SignIn/>} />
      <Route path="/signIn" element={<SignIn/>} />
      <Route path="/signUp" element={<SignUp/>} />
      <Route path="/reset" element={<UpdateDetails/>} />
      <Route path="/editProfile/:User" element={<ProfileSetup/>} />
      <Route path="/Profile/:User/:Page" element={<Profile/>} />
    </Routes>
  )
}

export default PublicRoute;