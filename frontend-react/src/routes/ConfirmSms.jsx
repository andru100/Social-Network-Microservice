import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import ConfirmEmail from './ConfirmEmail';

export default function ConfirmSms (props) {
	
	//const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage2");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [email, setEmail] = React.useState(props.email);
	const [mobile, setMobile] = React.useState(props.mobile);
	const [otp_sms, setOtp_sms] = React.useState("");

	async function verify(){
		
		const otp = document.getElementById('signupmobileotp').value;
			
		let signindata = {data:{
			Username: props.username,
			Password: props.password,
			Email:  props.email,
			Mobile: props.mobile,
			OTP_Mobile: otp,
			RequestType: "stage2"
			}
		}

		let gqlRequest = "mutation SignUp ($data: SecurityCheckInput!){ SignUp(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signup')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("Unable to send sms, please try again", response.errors[0].message )
			
		} else { // if password is a match redirect to profile page
			alert("Please confirm your email")
			setUsername(props.username)
			setPassword(props.password)
			setEmail(props.email)
			setMobile(props.mobile)
			setOtp_sms(otp)
			setStage("stage3")
		} 
		
		
	}

	return ( 
		<>
		{stage === "stage3" ?

      		<ConfirmEmail username={username} password={password} email={email} mobile={mobile} otp_sms={otp_sms} />
			:
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
					<input className="input100" type="text" name="mobileotp" id="signupmobileotp"/>
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
		}
		</> 
	)
};