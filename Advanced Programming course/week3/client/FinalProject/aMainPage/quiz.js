const quizContainer = document.getElementById('quiz');
const submitButton = document.getElementById('submitQuiz');
const resultContainer = document.getElementById('quizResult');
const restartButton = document.getElementById('restartQuiz');
const quizData = [
    {
        question: 'Which cuisine is known for sushi and sashimi?',
        options: ['Japanese', 'Italian', 'Mexican', 'Chinese'],
        answer: 'Japanese'
    },
    {
        question: 'What country is famous for its paella dish?',
        options: ['Spain', 'France', 'Italy', 'Greece'],
        answer: 'Spain'
    },
    {
        question: 'Which cuisine uses a lot of spices such as cumin, turmeric, and coriander?',
        options: ['Indian', 'Thai', 'Vietnamese', 'Greek'],
        answer: 'Indian'
    },
    {
        question: 'What is the main ingredient in hummus?',
        options: ['Chickpeas', 'Tomatoes', 'Eggplant', 'Pasta'],
        answer: 'Chickpeas'
    },
    {
        question: 'Which cuisine is known for its use of beans, tortillas, and chili peppers?',
        options: ['Mexican', 'Chinese', 'Japanese', 'Italian'],
        answer: 'Mexican'
    },
    {
        question: 'What is the traditional dish of the Thanksgiving meal in the United States?',
        options: ['Turkey', 'Pizza', 'Sushi', 'Pasta'],
        answer: 'Turkey'
    },
    {
        question: 'What type of food is traditionally eaten in the Japanese New Year?',
        options: ['Osechi', 'Burger', 'Pizza', 'Taco'],
        answer: 'Osechi'
    },
    {
        question: 'What is the key ingredient in a Margherita pizza?',
        options: ['Tomato', 'Bacon', 'Mushrooms', 'Pineapple'],
        answer: 'Tomato'
    },
    {
        question: 'What is the primary ingredient in a traditional Greek moussaka?',
        options: ['Eggplant', 'Zucchini', 'Potato', 'Carrot'],
        answer: 'Eggplant'
    },
    {
        question: 'Which cuisine is known for its use of kimchi and bulgogi?',
        options: ['Korean', 'Japanese', 'Italian', 'Thai'],
        answer: 'Korean'
    }
];

// Function to display the quiz questions
function displayQuiz() {
    const questions = document.querySelectorAll('.question');
    let index = 0;

    function showQuestion() {
        questions.forEach(question => {
            question.classList.remove('visible');
        });
        questions[index].classList.add('visible');
        index++;
        if (index === questions.length) {
            submitButton.style.display = 'block';
        }
    }

    showQuestion();

    submitButton.addEventListener('click', () => {
        if (index < questions.length) {
            showQuestion();
        } else {
            calculateResult();
        }
    });

    restartButton.addEventListener('click', () => {
        index = 0;
        createDropdownOptions();
        resultContainer.textContent = '';
        submitButton.style.display = 'block';
        showQuestion();
    });
}

// Function to shuffle array elements randomly
function shuffleArray(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}

// Function to calculate the quiz result
function calculateResult() {
    const answerContainers = quizContainer.querySelectorAll('.question select');
    let score = 0;

    quizData.forEach((data, index) => {
        const userAnswer = answerContainers[index].value;

        if (userAnswer === data.answer) {
            score++;
        }
    });

    displayResultMessage(score);
}

// Function to display result messages based on the score
function displayResultMessage(score) {
    let message = '';
    if (score >= 8) {
        message = "Congratulations, you know the cuisines of the world very well!";
    } else if (score >= 4 && score <= 7) {
        message = "Well done! You have a good knowledge of different cuisines.";
    } else {
        message = "Keep exploring! You can learn more about cuisines.";
    }

    resultContainer.textContent = message;
}

// Function to generate options with random correct answers for each question
function generateOptionsWithOptionsRandomized(data) {
    const correctAnswer = data.answer;
    const options = data.options.filter(option => option !== correctAnswer);
    options.push(correctAnswer);
    const randomizedOptions = shuffleArray(options);
    return randomizedOptions;
}
// Function to create the dropdown options for each question
function createDropdownOptions() {
    const dropdowns = document.querySelectorAll('.question select');
    quizData.forEach((data, index) => {
        const options = generateOptionsWithOptionsRandomized(data);
        dropdowns[index].innerHTML = options.map(option => `<option value="${option}">${option}</option>`).join('');
    });
}
// Restart quiz when the page loads
restartButton.click();

// Display the quiz on page load
displayQuiz();
createDropdownOptions();

