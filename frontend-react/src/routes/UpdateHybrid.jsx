import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import { useAlert } from "react-alert";
import RequestOTP from './RequestOTP';
import UpdateDetails from './UpdateDetails';
import Home from './Home';

export default function UpdateHybrid (props) {
	
	const Navigate = useNavigate();
	const [page, setPage] = React.useState("default");
	const [rendertype, setRenderType] = React.useState(props.rendertype); // an remove this and set to auto take session and get a proceed and address clue and use as a point on seurity.
	const [username, setUsername] = React.useState(props.username);
	const [authtype, setAuthType] = React.useState("init");
	const [address, setAddress] = React.useState("init");
	const [emailotp, setEmailOtp] = React.useState("");
	const [mobileotp, setMobileOtp] = React.useState("");
	const [updatetype, setUpdateType] = React.useState(props.updatetype);
	const [resetdata, setResetData] = React.useState({
		Username: props.username,
		UpdateType: props.updatetype,
		RequestType: "stage2",
	});
	const alert = useAlert()

	

	console.log("UpdateHybrid rendered: stage: always 2", "rendertype: ", rendertype, "authtype: ", authtype, "resetuserdata: ", resetdata)

		
	

	async function updateHybrid(){


		let updatedata = resetdata

		
		//console.log("authtype is: ", authtype)


		if (rendertype.includes("sms")) {
			const sms = document.getElementById('sms').value;
			updatedata.OTP_Mobile = sms
			
		}
		if (rendertype.includes("email")) {
			const email = document.getElementById('email').value;
			updatedata.OTP_Email = email
			
		}

		if (rendertype.includes("password")) {
			const password = document.getElementById('password').value;
			updatedata.Password = password
		}
		if (rendertype.includes("update")) {
			const updatevalue = document.getElementById('resethybrid').value;
			updatedata.UpdateData = updatevalue
			updatedata.RequestType = "update"
		}

		let gqlRequest = "mutation SecureUpdate ($data: SecurityCheckInput!){ SecureUpdate (input: $data) { Token AuthType MobClue EmailClue} }"
		
		let response = await SendData(gqlRequest, {data: updatedata}, 'secureupdate')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			console.log("error updating", response.errors[0].message )
			alertError(response.errors[0].message)
			return false
			
		} else {
			switch (response.data.SecureUpdate.Token){
				case "proceed":
					setResetData(updatedata)
					setAddress([response.data.SecureUpdate.MobClue, response.data.SecureUpdate.EmailClue])
					console.log("setting rendertype to next auth type from server", response.data.SecureUpdate.AuthType)
					setRenderType(response.data.SecureUpdate.AuthType)//update auth from server for next step/type of question dnd 
					break;
				case "update":
					setResetData(updatedata)
					console.log("server requesting update data ")
					setRenderType("update")//update auth from server for next step/type of question dnd 
					break;
				default:
					alert.show(updatetype + " updated")
					localStorage.setItem('jwt_token', response.data.SecureUpdate.Token) // Store JWT in storage
					if (updatedata.UpdateType === "Username") {
						updatedata.Username = updatedata.UpdateData
					}
					setResetData(updatedata)
					setPage("home")
					break;
			
			}
		
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
			//setAuthType("email")
			setRenderType("email")
			
		} else {
			ResendOTP("!sms")
			//setAuthType("sms")
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
				<button className="login100-form-btn" type="button" onClick={updateHybrid}>
						Verify
				</button>
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
				<div >
					<button  type="button" onClick={() => ChangeAuthType()}>
						dont have acccess to email
					</button>
				</div>
				<button className="login100-form-btn" type="button" onClick={updateHybrid}>
						Verify
				</button>
				
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
				<button className="login100-form-btn" type="button" onClick={updateHybrid}>
						Verify
				</button>
				
				
			</>
		)
	}

	function Update () {
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
					<button onClick={() => setPage("updatedetails")}>Back to account details</button>
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

	// const authTypeMap = {
	// 	"email": <ConfirmEmail/>,
	// 	"sms": <ConfirmSms/>,
	// 	"both": <HighSecurity/>,
	// 	"none": <ConfirmPassword/>
	// }


	function ConfirmUpdate () {
		return (
			<>	
				<Header/>
				{rendertype === "password" && <ConfirmPassword/>}
				{rendertype === "email" && <ConfirmEmail/>}
				{rendertype === "sms" && <ConfirmSms/>}
				{rendertype === "update" && <Update/>}
			</>
		)
	}

	return ( 
		<>
			
			{page === "default" && <ConfirmUpdate /> }
			{page === "updatedetails" && <UpdateDetails username={resetdata.Username}/>}
			{page === "home" && <Home sessionuser={resetdata.Username} viewing ={resetdata.Username} page={"home"}/>}
			
		</>
	)
};


// if password is a match redirect to profile page
			// switch (stage) {
			// 	case "stage1":
			// 		console.log("response is proceed in stage 1 setting states")
			// 		//setResetData({...resetdata, RequestType: "stage2"});
			// 		//setStage("stage2")
			// 		console.log("have set stage2, setting rendertype to authtype returned:", response.data.SecureUpdate.AuthType, "authtype piped in", authtype)
			// 		setRenderType(authtype)//update auth from default email everyone gets 
			// 		break;
			// 	case "stage2":
			// 		//setMobileOtp(response.data.OTP_Mobile)
			// 		setResetData({...resetdata, RequestType: "update"});
			// 		setStage("update")
			// 		setRenderType("update")
			// 		break;
			// 	case "update":
			// 		alert.show(updatetype + " updated")
			// 		localStorage.setItem('jwt_token', response.data.SecureUpdate.Token) // Store JWT in storage
			// 		resetdata.UpdateType === "Username" && setResetData({...resetdata,  Username: resetdata.UpdateData  });
					  
			// 		  //setUsername(updatevalue) 
			// 		setPage("home")
			// 		break;