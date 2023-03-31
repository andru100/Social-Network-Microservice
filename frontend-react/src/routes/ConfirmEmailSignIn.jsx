import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";
import ResetRequest from './UpdateDetails';
import RequestOTP from './RequestOTP';

export default function ConfirmEmailSignIn (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("email");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [address, setAddress] = React.useState(props.address)
	const [otp_sms, setOtp_sms] = React.useState("");
	const alert = useAlert()

	async function Verify(){

		console.log("password in emailconfurm:", password)
		
		const otp = document.getElementById('signinemailotp').value;
			
		let signindata = {data:{
			Username: username,
			Password: password,
			OTP_Email: otp,
			RequestType: "stage2"
			}
		}

		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
	
		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("user is not signed in", response.errors[0].message )
			return false
			
		  } else { // if password and email otp is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			Navigate ("/Profile/" + username + "/home")
		} 
		
		
		
	}

	async function resendOTP (requestType) {
		RequestOTP(username, requestType).then((response) => {
			if (( "errors" in response )) {
				console.log("error bak from otp request", response.errors[0].message)
				alertError(response.errors[0].message)
			}
		})
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
			<>
				<span className="login100-form-title p-b-53">
				<h3><i class="fa fa-lock fa-4x"></i></h3>
				</span>
				<div>
					
				</div>
				
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						{"Enter code sent to your email " + address[1] + "*******"}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="emailotp" id="signinemailotp"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<button  type="button" onClick={() => resendOTP("email")}>
						resend ccode
					</button>
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
				<div className="container-login100-form-btn m-t-17">
					<button className="login100-form-btn" type="button" onClick={() => Verify()}>
						Authenticate
					</button>
				</div>
				<div className="w-full text-center p-t-55">
					<button onClick={() => window.location.reload(false)}>Back to log in</button>
				</div>

				
			</>
			
		</>
	
)
};