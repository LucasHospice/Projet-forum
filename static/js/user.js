fetch('/users').then(function (response) {
    return response.json()
}).then(function (response) {
    const url = window.location.href.split('/')
    return response.filter(user => user.ID == url[url.length - 1])[0]
}).then(async response => {
    document.getElementById("profilname").innerText = response.Pseudo
    document.getElementById("profilrank").innerText = (response.Level == "1") ? "Utilisateur" : (response.Level == "2") ? "Modérateur" : "Administrateur"
    let pp = document.getElementById("profilavatar")
    pp.src = "../static/img/avatar/" + response.ProfilePic + ".png"
    const posts = await fetch("/posts").then(data => data.json()).then(data => data.filter(post => post.UserId == response.ID))
    posts.forEach(post => {
        let postsDiv = document.getElementById("profilpostedposts")
        let newDiv = document.createElement("div")
        newDiv.innerText = post.Title.String
        newDiv.style.cssText = "margin-top: 5%; cursor:pointer;border-top: 0.1rem solid; background-color:#ffffff;height:20%"
        newDiv.addEventListener("click", function () { location.href = "/topic/" + post.ID })
        postsDiv.appendChild(newDiv)
    })

    const upvotes = await fetch("/upvote").then(data => data.json()).then(async data => {
        let test = data.filter(upvote => {
            return upvote.UserId == response.ID
        })
        test = test.map(value => value.PostId)

        return await fetch("/posts").then(dat => dat.json()).then(poss => poss.filter(pos => {
            return test.includes(String(pos.ID))
        }))
    })


    upvotes.forEach(upvote => {
        let upvoteDiv = document.getElementById("profillikedposts")
        let newDiv = document.createElement("div")
        newDiv.innerText = upvote.Content
        // newDiv2.innerText = upvote.Title.String

        // newDiv.innerText = upvote.UserID.String
        newDiv.style.cssText = "margin-top: 5%; cursor:pointer;border-top: 0.1rem solid; background-color:#ffffff;height:20%"
        newDiv.addEventListener("click", function () { location.href = "/topic/" + upvote.ID })

        upvoteDiv.appendChild(newDiv)
    })
    const reponses = await fetch("/posts").then(data => data.json()).then(data => {
        let test = data.filter(reponse => reponse.UserID == reponse.UserID)
        return test
    })
    reponses.forEach(reponse => {
        // if (reponse.UserID == reponse.UserID)
        let reponseDiv = document.getElementById("profilansweredposts")
        let newDiv = document.createElement("div")
        // newDiv2.innerText = upvote.Title.String

        newDiv.innerText = reponse.Content
        newDiv.style.cssText = "margin-top: 5%; cursor:pointer;border-top: 0.1rem solid; background-color:#ffffff;height:20%"
        newDiv.addEventListener("click", function () { location.href = "/topic/" + reponse.ID })

        reponseDiv.appendChild(newDiv)
    })
})