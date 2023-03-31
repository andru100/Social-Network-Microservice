import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";
import ResetRequest from './UpdateDetails';
import RequestOTP from './RequestOTP';
import Home from './Home';

export default function SigninHybrid (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage2");
	const [page, setPage] = React.useState("default");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [authtype, setAuthType] = React.useState(props.authtype);
	const [address, setAddress] = React.useState(props.address);
	const [otp_sms, setOtp_sms] = React.useState("");
	const alert = useAlert()

	async function Verify(){

		console.log("password in smsconfurm:", password)
			
		let signindata = {data:{
			Username: username,
			Password: password,
			RequestType: "stage2"
			}
		}

		let getValues = []
		authtype === "both" ? getValues = ["sms", "email"] : authtype === "none" ? getValues = ["password"] : getValues = [authtype]
		
		console.log("getValues", getValues)
		if (getValues.includes("sms")) {
			const sms = document.getElementById('sms').value;
			signindata.data.OTP_Mobile = sms
		}
		if (getValues.includes("email")) {
			const email = document.getElementById('email').value;
			signindata.data.OTP_Email = email
		}
		



		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("user is not signed in", response.errors[0].message )
			return false
			
		  } else { // if password and Mobile otp is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			setPage("profile")
			//Navigate ("/Profile/" + username + "/home")
		} 
		
		
		
	}

	async function ResendOTP (requestType) {
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

	function ChangeAuthType () {

		console.log("change auth type called", authtype)
		//fun is only called when type is sms as button only shown on that page. added switch from email for future
		if (authtype === "sms") {
			ResendOTP("!email")
			setAuthType("email")
			
		} else {
			ResendOTP("!sms")
			setAuthType("sms")
			
		}
	}


	function ConfirmSms () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						{"Enter code sent to *******" + address[0]}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="mobileotp" id="sms"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<button  type="button" onClick={() => ResendOTP("sms")}>
						resend code
					</button>
				</div>
				<div >
					<button  type="button" onClick={() => ChangeAuthType()}>
						dont have acccess to phone
					</button>
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
			</>
		)
	}

	function ConfirmEmail () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						{"Enter code sent to email " + address[1] + "*******"}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="emailotp" id="email"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<button  type="button" onClick={() => ResendOTP("email")}>
						resend code
					</button>
				</div>
				<div >
					<button  type="button" onClick={() => ChangeAuthType()}>
						dont have acccess to email
					</button>
				</div>
				
			</>
		)
	}

	function SigninHybrid () {
		return (
			<>
				<div className="p-t-31 p-b-9">
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
				<div className="container-login100-form-btn m-t-17">
					<button className="login100-form-btn" type="button" onClick={Verify}>
						Sign in
					</button>
				</div>
				<div className="w-full text-center p-t-55">
					<button onClick={() => window.location.reload(false)}>Back to log in</button>
				</div>
			</>
		)
	}

	function Header () {
		return (
			<>
				<span className="login100-form-title p-b-53">
					<h3><i class="fa fa-lock fa-4x"></i></h3>
				</span>
				<div>
				</div>
			</>
		)
	}
	

	
	
	

	function HighSecurity () {
		return (
			<>
				<ConfirmSms/>
				<ConfirmEmail/>
			</>
		)
	}

	const authTypeMap = {
		"email": <ConfirmEmail/>,
		"sms": <ConfirmSms/>,
		"both": <HighSecurity/>,
	}

	function ConfirmSignin () {
		return (
			<>
				<Header/>
				{authTypeMap[authtype]}
				<SigninHybrid/>
			</>
		)
	}


	return ( 
		<>
			{page === "default" ? <ConfirmSignin/> : <Home/>}
			
		</>
)
};