/*======SHOW MENU======*/
const navMenu = document.getElementById("nav-menu"),
  navToggle = document.getElementById("nav-toggle"),
  navClose = document.getElementById("nav-close");

/*======MENU SHOW======*/
if (navToggle) {
  navToggle.addEventListener("click", () => {
    navMenu.classList.add("show-menu");
  });
}

/*======MENU HIDDEN=======*/
if (navClose) {
  navClose.addEventListener("click", () => {
    navMenu.classList.remove("show-menu");
  });
}

const navLink = document.querySelectorAll(".nav__link");

const linkAction = () => {
  const navMenu = document.getElementById("nav-menu");
  navMenu.classList.remove("show-menu");
};
navLink.forEach((n) => n.addEventListener("click", linkAction));

/* CHANGE BACKGROUND HEADER*/
const scrollHeader = () => {
  const header = document.getElementById("header");
  //when the scroll is greater than 50 viewport height, and the scroll header class to the header tag
  this.scrollY >= 50
    ? header.classList.add("bg-header")
    : header.classList.remove("bg-header");
};
window.addEventListener("scroll", scrollHeader);

/*SCROLL SECTIONS ACTIVE LINK*/
const sections = document.querySelectorAll("section[id]");

const scrollActive = () => {
  const scrollY = window.pageYOffset;

  sections.forEach((current) => {
    const sectionHeight = current.offsetHeight;
    const sectionTop = current.offsetTop - 58;
    const sectionId = current.getAttribute("id");
    const sectionClass = document.querySelector(
      '.nav__menu a[href*="' + sectionId + '"]'
    );

    if (scrollY > sectionTop && scrollY <= sectionTop + sectionHeight) {
      sectionClass.classList.add("active-link");
    } else {
      sectionClass.classList.remove("active-link");
    }
  });
};

window.addEventListener("scroll", scrollActive);

/*SHOW SCROLL UP*/
const scrollUp = () => {
  const scrollUp = document.getElementById("scroll-up");
  this.scrollY >= 350
    ? scrollUp.classList.add("show-scroll")
    : scrollUp.classList.remove("show-scroll");
};
window.addEventListener("scroll", scrollUp);

/*SCROLL REVEAL ANIMATION*/
const sr = ScrollReveal({
  origin: "top",
  distance: "60px",
  duration: 2500,
  delay: 400,
});

sr.reveal(`.home__data, .footer__container, .footer__group`);
sr.reveal(`.home__img`, { delay: 700, origin: "bottom" });
sr.reveal(`.logos__img, .program__card, .pricing__card`, { interval: 100 });
sr.reveal(`.choose__img, .calculate__content`, { origin: "left" });
sr.reveal(`.choose__content, .calculate__img`, { origin: "right" });

/*CALCULATE JS*/
const calculateForm = document.getElementById("calculate-form"),
  calculateCm = document.getElementById("calculate-cm"),
  calculateKg = document.getElementById("calculate-kg"),
  calculateMessage = document.getElementById("calculate-message");

const calculateBmi = (e) => {
  e.preventDefault();

  // Check if the fields have a value
  if (calculateCm.value === "" || calculateKg.value === "") {
    // Add and remove color
    calculateMessage.classList.remove("color-black");
    calculateMessage.classList.add("color-black");

    // Show message
    calculateMessage.textContent = "Fill in the Height and Weight 😡";

    // Remove message three seconds
    setTimeout(() => {
      calculateMessage.textContent = "";
    }, 4000);
  } else {
    //BMI FORMULA
    const cm = calculateCm.value / 100,
      kg = calculateKg.value,
      bmi = Math.round(kg / (cm * cm));

    //SHOW YOUR HEALTH STATUS
    if (bmi < 18.5) {
      //ADD COLOR AND DISPLAY MESSAGE
      calculateMessage.classList.add("color-black");
      calculateMessage.textContent = `Your BMI is ${bmi} and you are skinny 😔`;
    } else if (bmi < 25) {
      calculateMessage.classList.add("color-black");
      calculateMessage.textContent = `Your BMI is ${bmi} and you are healthy 🥳`;
    } else {
      calculateMessage.classList.add("color-black");
      calculateMessage.textContent = `Your BMI is ${bmi} and you are overweight 😔`;
    }

    //TO CLEAR THE INPUT FIELD
    calculateCm.value = "";
    calculateKg.value = "";

    //REMOVE MESSAGE FOUR SECONDS
    setTimeout(() => {
      calculateMessage.textContent = "";
    }, 4000);
  }
};

calculateForm?.addEventListener("submit", calculateBmi);

/*EMAIL JS*/
const contactForm = document.getElementById("contact-form"),
  contactMessage = document.getElementById("contact-message"),
  contactUser = document.getElementById("contact-user");

const sendEmail = (e) => {
  e.preventDefault();

  //CHECK IF THE FIELD HAS A VALUE
  if (contactUser.value === "") {
    //ADD AND REMOVE COLOR
    contactMessage.classList.remove("color-green");
    contactMessage.classList.add("color-red");

    //SHOW MESSAGE
    contactMessage.textContent = `You must enter your email 🤦`;

    //REMOVE MESSAGE THREE SECONDS
    setTimeout(() => {
      contactMessage.textContent = "";
    }, 3000);
  } else {
    //SERVICEID - TEMPLATEID - #FORM - PUBLICKEY
    emailjs
      .sendForm(
        "service_9ryh7td",
        "template_3argx64",
        "#contact-form",
        "kEeLyVneupuNvs2EO"
      )
      .then(
        () => {
          //SHOW MESSAGE AND ADD COLOR
          contactMessage.classList.add("color-green");
          contactMessage.textContent = "You registered successfully 💪";

          //REMOVE MESSAGE AFTER THREE SECONDS
          setTimeout(() => {
            contactMessage.textContent = "";
          }, 3000);
        },
        (error) => {
          //MAIL SENDING ERROR
          alert("OOPS! SOMETHING HAS FAILED...", error);
        }
      );
    //TO CLEAR THE INPUT FIELD
    contactUser.value = "";
  }
};

contactForm.addEventListener("submit", sendEmail);
