import React from 'react';
import './css/purified.css';
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