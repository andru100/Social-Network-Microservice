import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";
import RequestOTP from './RequestOTP';

export default function ResetPassword (props) {
	
	const Navigate = useNavigate();
	const [stage, setStage] = React.useState("stage2");
	const [username, setUsername] = React.useState(props.username);
	const [address, setAddress] = React.useState(props.address);
	const [otp_sms, setOtp_sms] = React.useState("");
	const alert = useAlert()

	

	async function updatepassword(){

		const mobileotp = document.getElementById('resetmobileotp').value;
		const emaileotp = document.getElementById('resetemaileotp').value;
		const newpassword = document.getElementById('resetnewpassword').value;
			
		let resetdata = {data: {
			Username: username,
			OTP_Email: emaileotp,
			OTP_Mobile: mobileotp,
			RequestType: "update",
			UpdateType: "Password",
			UpdateData: newpassword,
			}
		}

		let gqlRequest = "mutation SecureUpdate ($data: SecurityCheckInput!){ SecureUpdate (input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, resetdata, 'secureupdate')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("error resetting", response.errors[0].message )
			alertError(response.errors[0].message)
			return false
			
		  } else {// if password is a match redirect to profile page
			alert.show("password reset")
			localStorage.setItem('jwt_token', response.data.SecureUpdate.Token) // Store JWT in storage
			Navigate ("/Profile/" + username + "/home")
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

	return ( 
		<>
			<span className="login100-form-title p-b-53">
			<h3><i class="fa fa-lock fa-4x"></i></h3>
			</span>
			<div>
			</div>
			<div className="p-t-31 p-b-9">
				<span className="txt1">
					{"Enter code sent to *******" + address[0]}
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "verification code is required">
				<input className="input100" type="text" name="mobileotp" id="resetmobileotp"/>
				<span className="focus-input100"></span>
			</div>
			<div >
				<button  type="button" onClick={() => ResendOTP("sms")}>
					resend code
				</button>
			</div>
			<div className="container-login100-form-btn m-t-17">
			</div>
			<div className="p-t-31 p-b-9" >
				<span className="txt1" >
					{"Enter code sent to your email " + address[1] + "*******"}
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "verification code is required">
				<input className="input100" type="text" name="emaileotp" id="resetemaileotp"/>
				<span className="focus-input100"></span>
			</div>
			<div >
				<button  type="button" onClick={() => RequestOTP("email")}>
					resend ccode
				</button>
			</div>
			<div className="container-login100-form-btn m-t-17">
			</div>
			<div className="p-t-31 p-b-9">
				<span className="txt1">
					Enter your new password
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "verification code is required">
				<input className="input100" type="text" name="newpassword" id="resetnewpassword"/>
				<span className="focus-input100"></span>
			</div>
			<div className="container-login100-form-btn m-t-17">
			</div>
			<div className="container-login100-form-btn m-t-17">
				<button className="login100-form-btn" type="button" onClick={updatepassword}>
					Update
				</button>
			</div>
			<div className="w-full text-center p-t-55">
				<button onClick={() => window.location.reload(false)}>Back to log in</button>
			</div>

			
		</>
)
};