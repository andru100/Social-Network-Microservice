import React from 'react';
import './css/main.css';
import './css/util.css';
import './css/homePage.css';
import './css/displayPics.css';
import PublicRoutes from './PublicRoutes';
// import PrivateRoutes from './PrivateRoutes';

function App() {

  return (
  <div>
    <div className="limiter">
      <div className="container-login100" >
        <div className="wrap-login100 p-l-110 p-r-110 p-t-62 p-b-33">
          <form className="login100-form validate-form flex-sb flex-w">
            <PublicRoutes />
          </form>
        </div>
      </div>
    </div>
  </div>
  );
}

export default App