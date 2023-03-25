import React from 'react';
import { useEffect } from "react";
import { useNavigate, Link } from 'react-router-dom';
import ChkAuth from './chkAuth';
import Google from '../images/icons/icon-google.png'
import SendData from './SendData';
import ResetPassword from './ResetPassword';

export default function RenderResetPassword () {

	const Navigate = useNavigate();
	const [updateType, setUpdateType] = React.useState("init");

	useEffect( () => { //check if signed in and go to profile page
		ChkAuth().then(user => {
		  if (user) {
		  Navigate("/profile/" + user + "/home")
		  } 
		})
	},[]);

	
	

	async function RequestOTP (requestType) {
			
		let signindata = {data: {
			Username: username,
			RequestType: requestType
			}
		}

		let gqlRequest = "mutation SecureUpdate ($data: SecurityCheckInput!){ SecureUpdate(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'secureupdate')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			console.log("Unable to re-send OTP", response.errors[0].message )
			alertError(response.errors[0].message)
			
		} else { // if password is a match redirect to profile page
			alert.show("We have re-sent your code")
			//setStage("stage2")
			
		} 

	}
	
  return (  
    <>
	{updateType === "Password" ?

	<ConfirmSmsSignIn username={username} password={password} />
	:
	<>
			<span className="login100-form-title p-b-53">
			<h3><i class="fa fa-lock fa-4x"></i></h3>
			</span>
			<div className="container-login100-form-btn m-t-17">
				<button className="login100-form-btn" type="button" onClick={verify}>
					Authenticate
				</button>
			</div>
			<div className="w-full text-center p-t-55">
				<button onClick={() => window.location.reload(false)}>Back to log in</button>
			</div>

			
		</>

	}
	</>


  )
};