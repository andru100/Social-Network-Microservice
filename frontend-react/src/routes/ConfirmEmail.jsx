import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';

export default function ConfirmEmail (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage3");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [email, setEmail] = React.useState(props.email);
	const [mobile, setMobile] = React.useState(props.mobile);
	const [otp_sms, setOtp_sms] = React.useState(props.otp_sms);

	async function verify(){
		
		const otp = document.getElementById('signupemailotp').value;
			
		let signindata = {data:{
			Username: props.username,
			Password: props.password,
			Email:  props.email,
			Mobile: props.mobile,
			OTP_Mobile: props.otp_sms,
			OTP_Email: otp,
			RequestType: stage
			}
		}

		let gqlRequest = "mutation SignUp ($data: SecurityCheckInput!){ SignUp(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signup')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("Unable to send email", response.errors[0].message )
			
		} else {// if password is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignUp.Token) // Store JWT in storage
      		alert("Welcome to the club! Please setup your profile")
			Navigate ("/editProfile/"+username)
		} 
		
		
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
						Enter the code sent to {props.email}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="emailotp" id="signupemailotp"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<Link to="/reset">Resend code</Link>
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
				<div className="container-login100-form-btn m-t-17">
					<button className="login100-form-btn" type="button" onClick={verify}>
						Authenticate
					</button>
				</div>
				<div className="w-full text-center p-t-55">
					<Link to="/signIn">Back to log in</Link>
				</div>
				
			</>
	)
};