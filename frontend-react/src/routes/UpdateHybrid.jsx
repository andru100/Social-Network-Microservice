import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";
import RequestOTP from './RequestOTP';

export default function UpdateHybrid (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage2");
	const [username, setUsername] = React.useState(props.username);
	const [authtype, setAuthType] = React.useState(props.authtype);
	const [address, setAddress] = React.useState(props.address);
	const [updatetype, setUpdateType] = React.useState(props.updatetype);
	const [otp_sms, setOtp_sms] = React.useState("");
	const alert = useAlert()

	console.log("UpdateHybrid username: Password:", username)
	

	async function updateHybrid(){

		//const emailotp = document.getElementById('resetmobileotp').value;
		const updatevalue = document.getElementById('resethybrid').value;
			
		let resetdata = {data: {
			Username: username,
			RequestType: "update",
			UpdateType: updatetype,
			UpdateData: updatevalue,
			}
		}

		let getValues = []
		authtype === "both" ? getValues = ["sms", "email"] : authtype === "none" ? getValues = ["password"] : getValues = [authtype]
		
		console.log("getValues", getValues)
		if (getValues.includes("sms")) {
			const sms = document.getElementById('sms').value;
			resetdata.data.OTP_Mobile = sms
		}
		if (getValues.includes("email")) {
			const email = document.getElementById('email').value;
			resetdata.data.OTP_Email = email
		}

		if (getValues.includes("password")) {
			const password = document.getElementById('password').value;
			resetdata.data.Password = password
		}

		//if (updatetype !== "Password") {
			const jwt = localStorage.getItem('jwt_token');
			resetdata.data.Token = jwt
		//}

		let gqlRequest = "mutation SecureUpdate ($data: SecurityCheckInput!){ SecureUpdate (input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, resetdata, 'secureupdate')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("error updating", response.errors[0].message )
			alertError(response.errors[0].message)
			return false
			
		  } else {// if password is a match redirect to profile page
			alert.show(updatetype + " updated")
			localStorage.setItem('jwt_token', response.data.SecureUpdate.Token) // Store JWT in storage
			//updatevalue === "Username" && props.setUsername(updatevalue) 
			Navigate('/')
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

	function ConfirmPassword () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						Enter your password
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "password is required">
					<input className="input100" type="text" name="pass" id="password"/>
					<span className="focus-input100"></span>
				</div>
				
			</>
		)
	}

	function UpdateHybrid () {
		return (
			<>
				<div className="p-t-31 p-b-9">
				<span className="txt1">
					{"Enter your new " + updatetype}
				</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="reset" id="resethybrid"/>
					<span className="focus-input100"></span>
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
				<div className="container-login100-form-btn m-t-17">
					<button className="login100-form-btn" type="button" onClick={updateHybrid}>
						Update
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
		"none": <ConfirmPassword/>
	}


	return ( 
		<>
			<Header/>
			{authTypeMap[authtype]}
			<UpdateHybrid/>
		</>
)
};