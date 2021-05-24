let dropArea = document.getElementById('drop-area')

;
['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
    dropArea.addEventListener(eventName, preventDefaults, false)
})

;
['dragenter', 'dragover'].forEach(eventName => {
    dropArea.addEventListener(eventName, highlight, false)
})

;
['dragleave', 'drop'].forEach(eventName => {
    dropArea.addEventListener(eventName, unhighlight, false)
})

function highlight(e) {
    dropArea.classList.add('highlight')
}

function unhighlight(e) {
    dropArea.classList.remove('highlight')
}

dropArea.addEventListener('drop', handleDrop, false)

function handleDrop(e) {
    let dt = e.dataTransfer
    let files = dt.files

    handleFiles(files)
}

function handleFiles(files) {
    files = [...files]
    files.forEach(shouldNotExceedFirestoreFieldSizeLimit)
    files.forEach(uploadFile)
}

function shouldNotExceedFirestoreFieldSizeLimit(file) {
    if (file.size > 900000) { // check max firestore field size limit: 1048487 minus prefix 'data:image/png;base64,'
        alert('Uploaded image size for image ' + file.name + ' is too large: should not exceed 900 KB');
    }
}

function uploadFile(file) {
    let url = '/api/uploads'
    let formData = new FormData()

    formData.append('photos', file)

    fetch(url, {
            method: 'POST',
            body: formData
        })
        .then((res) => res.json())
        .then((json) => json.forEach((res) => previewFile(file, res)))
        .catch((e) => { alert(e) })
}

function previewFile(file, res) {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onloadend = function() {
        let place = res.Landmarks[0]
        if (place == undefined) {
            place = "Unknown"
        }

        const img = createCard(reader.result, place)
        document.getElementById('gallery').appendChild(img)
    }
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