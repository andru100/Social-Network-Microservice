import React from 'react';
import { Routes, Route} from 'react-router-dom';
import SignUp from "./routes/SignUp";
import SignIn from "./routes/SignIn";
import ProfileSetup from "./routes/ProfileSetup";
import Profile from "./routes/Profile";

//Allows access to pages that have no authorisation and allows them to sign in/up

const PublicRoute = () => {
  return (
    <Routes>
      <Route path="/" element={<SignIn/>} />
      <Route path="/signIn" element={<SignIn/>} />
      <Route path="/signUp" element={<SignUp/>} />
      <Route path="/editProfile/:user" element={<ProfileSetup/>} />
      <Route path="/Profile/:Page/:user" element={<Profile/>} />
    </Routes>
  )
}

export default PublicRoute;