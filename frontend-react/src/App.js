import React from 'react';
import './css/main.css';
import './css/util.css';
import './css/homePage.css';
import './css/displayPics.css';
import 'bootstrap/dist/css/bootstrap.min.css'
import ThemeProvider from 'react-bootstrap/ThemeProvider'
import SignIn from './routes/SignIn';

function App() {

  return (
  <div>
    <SignIn/>
  </div>
  );
}

export default App