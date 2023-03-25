import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";

export default function ConfirmSmsSignIn (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage2");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [otp_sms, setOtp_sms] = React.useState("");
	const alert = useAlert()

	async function verify(){
		
		const otp = document.getElementById('signinmobileotp').value;
			
		let signindata = {data:{
			Username: props.username,
			Password: props.password,
			OTP_Mobile: otp,
			RequestType: "stage2"
			}
		}

		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
	
		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("user is not signed in", response.errors[0].message )
			return false
			
		  } else { // if password and mobile otp is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			Navigate ("/Profile/" + username + "/home")
		} 
		
		
		
	}

	async function RequestOTP () {
			
		let signindata = {data: {
			Username: username,
			Password: password,
			RequestType: "stage1"
			}
		}

		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			console.log("Unable to re-send OTP", response.errors[0].message )
			alertError(response.errors[0].message)
			
		} else { // if password is a match redirect to profile page
			alert.show("We have re-sent your code")
			//setStage("stage2")
			
		} 

	}

	
	function alertError(error){
		const delimiter = '= '
		const start = 2,
		tokens = error.split(delimiter).slice(start),
		result = tokens.join(delimiter); // those.that
		alert.show(result)
	}

	return ( 
		<>
			<span className="login100-form-title p-b-53">
			<h3><i class="fa fa-lock fa-4x"></i></h3>
			</span>
			<div>
				
			</div>
			
			<div className="p-t-31 p-b-9">
				<span className="txt1">
					Enter the code sent to {props.mobile}
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "verification code is required">
				<input className="input100" type="text" name="mobileotp" id="signinmobileotp"/>
				<span className="focus-input100"></span>
			</div>
			<div >
				<button  type="button" onClick={RequestOTP}>
					resend ccode
				</button>
			</div>
			<div className="container-login100-form-btn m-t-17">
			</div>
			<div className="container-login100-form-btn m-t-17">
				<button className="login100-form-btn" type="button" onClick={verify}>
					Authenticate
				</button>
			</div>
			<div className="w-full text-center p-t-55">
				<button onClick={() => window.location.reload(false)}>Back to log in</button>
			</div>

			
		</>
)
};