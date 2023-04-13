import React from 'react';
import { useAlert } from "react-alert";
import SendData from './SendData';
import RequestOTP from './RequestOTP';
import Home from './Home';

export default function SigninHybrid (props) {

	const [page, setPage] = React.useState("default");
	const [username, setUsername] = React.useState(props.username);
	const [password, setPassword] = React.useState(props.password);
	const [authtype, setAuthType] = React.useState(props.authtype);
	const [address, setAddress] = React.useState(props.address);
	const [rendertype, setRenderType] = React.useState(props.authtype);
	const [userdata, setUserdata] = React.useState({
		Username: props.username,
		Password: props.password,
		RequestType: "stage1"
	});

	const alert = useAlert()

	async function Verify(){

			
		let signindata = userdata
		
		if (rendertype === "sms") {
			const sms = document.getElementById('sms').value;
			signindata.OTP_Mobile = sms
		}
		if (rendertype === "email") {
			const email = document.getElementById('email').value;
			signindata.OTP_Email = email
		}
		
		if (rendertype === "oauth") {
			const passcode = document.getElementById('oauth').value;
			signindata.Oauth = passcode
		}



		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, {data:signindata}, 'signin')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("user is not signed in", response.errors[0].message )
			return false
			
		} 
		  
		if (response.data.SignIn.Token === "proceed") {
			//stage === "stage2" && setEmailOtp(signupdata.data.OTP_Email)
			setUserdata(signindata);
			//setStage("stage3")
			console.log("setting rendertype to next auth type from server", response.data.SignIn.AuthType)
			setRenderType(response.data.SignIn.AuthType)//update auth from server for next step/type of question dnd 
			

		} else {
			    console.log("user is signed in token is", response.data.SignIn.Token )
				localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			   	setPage("home")
		}
		
		
		
		
	}

	async function ResendOTP (requestType) {
		RequestOTP(username, requestType, "user").then((response) => {
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
		//func is only called when type is sms as button only shown on that page. added switch from email for future
		if (rendertype === "sms") {
			ResendOTP("!email")
			setRenderType("email")
			
		} else {
			ResendOTP("!sms")
			setRenderType("sms")
			
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
					<button  type="button" onClick={() => ResendOTP("!sms")}>
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
					<button  type="button" onClick={() => ResendOTP("!email")}>
						resend code
					</button>
				</div>
				{/* <div >
					<button  type="button" onClick={() => ChangeAuthType()}>
						dont have acccess to email
					</button>
				</div> */}
				
			</>
		)
	}

	function Authenticate () {
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
					<button onClick={() => setPage("default")}>Back to account details</button>
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
	

	
	
	

	// function HighSecurity () {
	// 	return (
	// 		<>
	// 			<ConfirmSms/>
	// 			<ConfirmEmail/>
	// 		</>
	// 	)
	// }

	// const authTypeMap = {
	// 	"email": <ConfirmEmail/>,
	// 	"sms": <ConfirmSms/>,
	// 	"both": <HighSecurity/>,
	// }

	function ConfirmSignin () {
		return (
			<>	
				<div className="limiter">
      				<div className="container-login100" >
       					<div className="wrap-login100 p-l-110 p-r-110 p-t-62 p-b-33">
          					<form className="login100-form validate-form flex-sb flex-w">
									<Header/>
									{rendertype === "email" && <ConfirmEmail/>}
									{rendertype === "sms" && <ConfirmSms/>}
									<Authenticate/>
							</form>
        				</div>
      				</div>
    			</div>
			</>
		)
	}


	return ( 
		<>
			{page === "default" && <ConfirmSignin/>}
			{page === "home" && <Home sessionuser={username} page={"home"} viewing={username}/>}
			
		</>
	)
};