function croissant(value) {
    innerHTML = ""
    if (value === "upvote") {
        fetch("/topics").then(function (response) {
            return response.json()
        }).then(response => filtre(response.sort((greater, lower) => greater.UpVote - lower.UpVote)))
        // } else if (value === "date") {
        //     //la function marche je n'arrive pas a l'afficher mais le tri ce fait
        //     //il faudrait je pense reussir a afficher les post par id car l'array est seulement rempli de Date et des ID des post
        //     fetch("/topics").then(function (response) {
        //         return response.json()
        //             .then(response => {
        //                 let allArrays = [];
        //                 for (let i = 0; i < response.length; i++) {
        //                     allArrays.push(new Date(response[i].Date))
        //                 }
        //                 allArrays.sort((a, b) => a - b)
        //                 filtre(allArrays)
        //             })
        //     })
    }
}

function decroissant(value) {
    innerHTML = ""
    if (value === "upvote") {
        fetch("/topics").then(function (response) {
            return response.json()
        }).then(response => filtre(response.sort((greater, lower) => lower.UpVote - greater.UpVote)))
        // } else if (value === "date") {
        //     //la function marche je n'arrive pas a l'afficher mais le tri ce fait 
        //     //il faudrait je pense reussir a afficher les post par id car l'array est seulement rempli de Date et des ID des post
        //     fetch("/topics").then(function (response) {
        //         return response.json()
        //             .then(response => {
        //                 let allArrays = [];
        //                 for (let i = 0; i < response.length; i++) {
        //                     allArrays.push([new Date(response[i].Date), response[i].ID])
        //                 }
        //                 allArrays.sort((a, b) => b[0] - a[0])
        //             })

        //     })
    }
}



// function filtreCategories(value) {
//     // mauvaise techno je pense les categories sont des ID a revoir entierement ou non 
//     innerHTML = ""
//     if (value === "feuilleton") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "feuillton")))
//     } else if (value === "jardin") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "jardin")))
//     } else if (value === "arnaque") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "arnaque")))
//     } else if (value === "education") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "education")))
//     } else if (value === "sante") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category == 0)))
//     } else if (value === "politique") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "politique")))
//     } else if (value === "nostalgie") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => console.log(topics.Category) == 1)))
//     } else if (value === "complot") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "complot")))
//     } else if (value === "ragot") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "ragot")))
//     } else if (value === "sport") {
//         fetch("/topics").then(function (response) {
//             return response.json()
//         }).then(response => filtre(response.filter((topics) => topics.Category === "sport")))
//     }
// }

function reinitialiser() {
    window.location.reload()
}


function filtre(response) {
    console.log(response + "|||||func filtre")
    let list = document.getElementById("topicList")
    for (let i in response) {
        const p = document.createElement('div')
        fetch('./static/components/postCard.txt')
            .then(response => response.text())
            .then(data => {
                console.log(response[i].Date)
                data = data.split("{Pseudo}").join(response[i].Title.String).split("{Content}").join(response[i].Content).split("{Date}").join(response[i].Date).split("{PostId}").join(response[i].ID).split("{UserId}").join(response[i].UserId).split("{UpVote}").join(response[i].UpVote).split("{PostId}").join(response[i].ID)
                // console.log("comp1 " + data)
                fetch("/categories").then(catt => catt.json()).then(function (catt) {
                    let category = catt.filter(obj => obj.ID == response[i].Category.Int64)[0]
                    data = data.split("{CatColor}").join(category.Color).split("{Category}").join(category.Name)
                    p.innerHTML = data
                    list.appendChild(p)
                })
            })
    }
}

