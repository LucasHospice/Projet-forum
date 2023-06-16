// export const displayPosts = () => {
fetch("/posts").then(function (response) {
    return response.json()
}).then(response => response.filter(post => post.IsTopic == 1)).then(function (response) {
    let list = document.getElementById("topicList")
    for (let i in response) {
        const p = document.createElement('div')
        fetch('./static/components/postCard.txt')
            .then(resp => resp.text())
            .then(data => {
                // Do something with your data
                data = data.split("{Pseudo}").join(response[i].Title.String).split("{Content}").join(response[i].Content).split("{Date}").join(response[i].Date).split("{PostId}").join(response[i].ID).split("{UserId}").join(response[i].UserId).split("{UpVote}").join(response[i].UpVote).split("{PostId}").join(response[i].ID)
                fetch("/categories").then(catt => catt.json()).then(function (catt) {
                    let category = catt.filter(obj => obj.ID == response[i].Category.Int64)[0]
                    data = data.split("{CatColor}").join(category.Color).split("{Category}").join(category.Name)
                    fetch("/users").then(users => users.json()).then(users => users.filter(user => user.ID == response[i].UserId)[0]).then(user => {
                        data = data.split("{ProfilePic}").join(user.ProfilePic)
                        p.innerHTML = data
                        list.appendChild(p)
                        fetch("/posts")
                    })
                })
            })
    }
}).catch(function (err) {
    console.log(err)
})