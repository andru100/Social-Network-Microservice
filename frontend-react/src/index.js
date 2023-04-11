import React from 'react';
import { render } from "react-dom";
import './index.css';
import App from './App';
import { transitions, positions, Provider as AlertProvider } from 'react-alert'
import AlertTemplate from 'react-alert-template-mui'

// optional configuration
const options = {
  // you can also just use 'bottom center'
  position: positions.BOTTOM_CENTER,
  timeout: 5000,
  offset: '30px',
  // you can also just use 'scale'
  transition: transitions.SCALE
}


const rootElement = document.getElementById("root");
render(

    <AlertProvider template={AlertTemplate} {...options}>
        <App />
    </AlertProvider>,
   rootElement
);

