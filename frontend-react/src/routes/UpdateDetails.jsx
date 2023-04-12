import React from 'react';
import { useEffect } from "react";
import RequestOTP from './RequestOTP';
import UpdateHybrid from './UpdateHybrid';
import Home from './Home';
import { useAlert } from "react-alert";

export default function RenderUpdateDetails (props) {

	const [authtype, setAuthType] = React.useState("init");
	const [username, setUsername] = React.useState(props.username)
	const [address, setAddress] = React.useState("unassigned")
	const [userDetails, setUserDetails] = React.useState({username: "", password: "", email: "", mobile: ""})
	const [name, setName] = React.useState("John Doe");
	const [email, setEmail] = React.useState("johndoe@gmail.com");
	const [phone, setPhone] = React.useState("555-555-5555");
	const [page, setPage] = React.useState("default");
	
  


	console.log("update details render start, username:", username, "page:", page)

	// async function GetOTP (requestType, pagetype) {
	// 	if (requestType == "none") {
	// 		setAuthType(requestType)
	// 		setPage(pagetype)
	// 		return
	// 	}

	// 	RequestOTP(username, requestType, "user").then((response) => {
	// 		if (( "errors" in response )) {
	// 			console.log("error bak from otp request", response.errors[0].message)
	// 			alertError(response.errors[0].message)
	// 		} else {
	// 			setAddress([response.data.SecureUpdate.MobClue, response.data.SecureUpdate.EmailClue])
	// 			setAuthType(response.data.SecureUpdate.AuthType)
	// 			console.log("otp response setting page:", pagetype)
	// 			setPage(pagetype)
				

	// 		}
	// 	})
	// }

	function StartUpdate (updatetype) {
		setPage(updatetype)
	}

	const handleUpdateName = (event) => {
		setName(event.target.value);
	};
	
	const handleUpdateEmail = (event) => {
		setEmail(event.target.value);
	};
	
	const handleUpdatePhone = (event) => {
		setPhone(event.target.value);
	};
	
	const handleUpdateAddress = (event) => {
		setAddress(event.target.value);
	};

	function alertError(error){
		const delimiter = '= '
		const start = 2
		const tokens = error.split(delimiter).slice(start)
		const result = tokens.join(delimiter); // those.that
		alert.show(result)
	}

	function GoHome () {
		console.log("going home, username:", username)
		setPage("home")
	}

	function LandingPage () {
		return (
			<>
				<div className="limiter">
      				<div className="container-login100" >
       					<div className="wrap-login100 p-l-110 p-r-110 p-t-62 p-b-33">
          					<form className="login100-form validate-form flex-sb flex-w">
								<div>
									<h1>Account Details</h1>
									<div>
										<label>Username:</label>
										<input type="text" value={name} onChange={handleUpdateName} />
										<button className="login100-form-btn" type="button" onClick={() => StartUpdate("Username")}>
											update username
										</button>
									</div>
									<div>
										<label>Email:</label>
										<input type="text" value={email} onChange={handleUpdateEmail} />
										<button className="login100-form-btn" type="button" onClick={() => StartUpdate("Email")}>
											update email
										</button>
									</div>
									<div>
										<label>Phone:</label>
										<input type="text" value={phone} onChange={handleUpdatePhone} />
										<button className="login100-form-btn" type="button" onClick={() => StartUpdate("Mobile")}>
											update mobile
										</button>
									</div>
									<div>
										<label>Password:</label>
										<input type="text" value={address} onChange={handleUpdateAddress} />
										<button className="login100-form-btn" type="button" onClick={() => StartUpdate("Password")}>
											update password
										</button>
									</div>
									<div>
										<button type="button" onClick={() => GoHome()}>
											back to home
										</button>
									</div>
								</div>
							</form>
						</div>
					</div>
				</div>
			</>
		)
	}


  return ( 
	<>
		
		{page === "default" ? <LandingPage /> : page === "home" ? <Home sessionuser={username} viewing= {username} page={"home"}/> : <UpdateHybrid username={username} /*  address ={address} */ updatetype = {page} rendertype={"password"}/* authtype ={authtype} *//>}
		
	</>
	)
};


