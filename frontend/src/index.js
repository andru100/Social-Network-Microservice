import React from 'react';
import { render } from "react-dom";
import { BrowserRouter, Routes, Route} from "react-router-dom";
import './index.css';
import App from './App';



const rootElement = document.getElementById("root");
render(
   <BrowserRouter>
    <Routes>
      <Route path="*" element={<App />}/>
    </Routes>
   </BrowserRouter>,
   rootElement
);

