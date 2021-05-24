const gallery = document.getElementById('gallery')
let cards = []
let creationDates = []

refresh2()
setInterval(refresh2, 5000)


function refresh() {
    fetch("/api/uploads?limit=10&offset=0", {
            method: 'GET',
        })
        .then((res) => res.json())
        .then((json) => json.forEach((res) => {
            previewFile(res)
        }))
        .catch((e) => { alert(e) })
}

function refresh2() {
    let lastCreationDate = creationDates.sort().reverse()[0] || '1970-01-01T00:00:00.000Z'
    fetch("/api/uploads?creationDate=" + lastCreationDate, {
            method: 'GET',
        })
        .then((res) => res.json())
        .then((json) => json.forEach((res) => {
            previewFile(res)
        }))
        .catch((e) => { alert(e) })
}

function previewFile(res) {
    let place = "Unknown"
    if (res.Landmarks !== undefined) {
        place = res.Landmarks[0]
    }

    const img = createCard(res.Image, place)
    if (!containsFile(res.Filename)) {
        document.getElementById('gallery').appendChild(img)
        cards.push(res.Filename)
        creationDates.push(res.CreationDate)
    }
}


function containsFile(filename) {
    for (i = 0; i < cards.length; i++) {
        if (cards[i] === filename) {
            return true
        }
    }
    return false
}

function createCard(data, landmarks) {
    const card = document.createElement('div')
    card.classList.add("card", "fade-in")

    const img = document.createElement('img')
    img.src = data

    const imgContainer = document.createElement('div')
    imgContainer.classList.add("img-container")

    const h4 = document.createElement('h4')

    const b = document.createElement('b')
    b.textContent = landmarks

    card.appendChild(img)
    card.appendChild(imgContainer)
    imgContainer.appendChild(h4)
    h4.appendChild(b)

    return card
}

function preventDefaults(e) {
    e.preventDefault()
    e.stopPropagation()
}