let monCompte = document.getElementById("monCompte");
let posts = document.getElementById("posts");
let paramètres = document.getElementById("paramètes");
let déconnexion = document.getElementById("déconnexion");
let photo = document.getElementById("photo");
let pseudo = document.getElementById("pseudo");
let mail = document.getElementById("e-mail");
let rang = document.getElementById("rang");
let liste = document.getElementById("liste");
let listePosts = document.getElementsByClassName("listePosts");
let déco = document.getElementById("déco");
let mesCommentaires = document.getElementById("mesCommentaires");
let commentaires = document.getElementById("commentaires");
let badges = document.getElementById("badges");
let déroul = document.getElementById("scroll");
let non = document.getElementById("non")
let infoDiv = document.getElementById("info")

monCompte.addEventListener("click", () => {


        infoDiv.style.backgroundColor = "#1C1F4C"
        photo.style.display = "block";
        pseudo.style.display = "block";
        mail.style.display = "block";
        rang.style.display = "block"
        badges.style.display = "block"
        liste.style.display = "none"
        déco.style.display = "none";
        commentaires.style.display = "none"
        déroul.style.display = "none"
        openPopup.style.display = "none"

}
)

posts.addEventListener("click", () => {


        déroul.style.display = "block"
        photo.style.display = "none";
        pseudo.style.display = "none";
        mail.style.display = "none";
        rang.style.display = "none";
        déco.style.display = "none";
        commentaires.style.display = "none"
        badges.style.display = "none"
        openPopup.style.display = "none"
        infoDiv.style.backgroundColor = "#698EA2"


}
)

déconnexion.addEventListener("click", () => {
        openPopup.style.display = "block"

})

non.addEventListener("click", () => {
        openPopup.style.display = "none"
})

function close() {
        console.log("test")
}

fetch("/cookies-data").then((response) => response.json()).then(data => {
        document.getElementById("PseudoDisplay").innerText = data.pseudo
        fetch("/users").then(response => response.json()).then(response => {
                let userData = response.filter(user => user.ID == data.user_id)[0]
                document.getElementById("EmailDisplay").innerText = userData.Mail
                document.getElementById("NumberDisplay").innerText = userData.Number
                document.getElementById("photo").src = "../static/img/avatar/" + userData.ProfilePic + ".png"
        })
})

fetch('/cookies-data').then(function (response) {
        return response.json()
}).then(async function (response) {
        return await fetch("/users").then(rush => rush.json()).then(rush => rush.filter(obj => obj.ID == response.user_id)[0])
}).then(async response => {

        const posts = await fetch("/posts").then(data => data.json()).then(data => data.filter(post => post.UserId == response.ID))
        posts.forEach(post => {
                let postsDiv = document.getElementById("postUpvotes")
                let newDiv = document.createElement("div")
                newDiv.innerText = post.Title.String
                newDiv.style.cssText = "margin-top: 5%; cursor:pointer;border-top: 0.1rem solid; background-color:#ffffff;height:20%"
                newDiv.addEventListener("click", function () { location.href = "/topic/" + post.ID })
                postsDiv.appendChild(newDiv)
        })
})