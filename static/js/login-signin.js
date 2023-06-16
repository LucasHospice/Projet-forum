const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');
const logBlocker = document.getElementById('logblocker')
const signIn = document.getElementById('menusignin')
const signUp = document.getElementById('menusignup')
const notLoggedSignIn = document.getElementById('notLoggedLogIn')
const notLoggedSignUp = document.getElementById('notLoggedLogUp')
const logError = document.getElementById('logerror')

signUpButton.addEventListener('click', () => {
    container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
    container.classList.remove("right-panel-active");
});

logBlocker.addEventListener('click', () => {
    document.getElementById('container').style.display = 'none'
    logBlocker.style.display = 'none'
    logError.style.display = 'none'
    document.body.style.overflow = 'auto'
});

signUp.addEventListener('click', () => {
    container.style.display = 'block'
    logBlocker.style.display = 'block'
    container.classList.add("right-panel-active");
    document.body.style.overflow = 'hidden'
});

signIn.addEventListener('click', () => {
    container.style.display = 'block'
    logBlocker.style.display = 'block'
    container.classList.remove("right-panel-active");
    document.body.style.overflow = 'hidden'
});

notLoggedSignUp.addEventListener('click', () => {
    container.style.display = 'block'
    logBlocker.style.display = 'block'
    container.classList.add("right-panel-active");
    document.body.style.overflow = 'hidden'
});

notLoggedSignIn.addEventListener('click', () => {
    container.style.display = 'block'
    logBlocker.style.display = 'block'
    container.classList.remove("right-panel-active");
    document.body.style.overflow = 'hidden'
});

const checkLoginForm = () => {
    fetch("/users").then(function (response) {
        return response.json()
    }).then(function (response) {
        for (let i in response) {
            let email = document.getElementById('email1')
            let mdp = document.getElementById('mdp1')
            if (response[i].Pseudo == email.value && response[i].Password == mdp.value) {
                document.getElementById('subscribe').submit()
                return
            }
            errorP = document.getElementById("loginError")
            errorP.innerText = "Identifiants incorrects"
        }
    }
    ).catch(function (err) {
        console.log(err)
    })
}