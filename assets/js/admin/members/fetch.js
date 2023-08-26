/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./views/admin/members/fetch.js":
/*!**************************************!*\
  !*** ./views/admin/members/fetch.js ***!
  \**************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   cancelSubscription: () => (/* binding */ cancelSubscription),\n/* harmony export */   fetchClients: () => (/* binding */ fetchClients),\n/* harmony export */   fetchMembers: () => (/* binding */ fetchMembers),\n/* harmony export */   fetchMembershipPlans: () => (/* binding */ fetchMembershipPlans),\n/* harmony export */   subscribe: () => (/* binding */ subscribe)\n/* harmony export */ });\nconst fetchMembers = async () => {\n  try {\n    const response = await fetch(\"/members\", {\n      headers: new Headers({ \"content-type\": \"application/json\" }),\n    });\n    const { data } = await response.json();\n    return data?.members ?? [];\n  } catch (error) {\n    console.error(error);\n    return;\n  }\n};\nconst fetchClients = async () => {\n  try {\n    const response = await fetch(\"/clients?type=unsubscribed\", {\n      headers: new Headers({ \"content-type\": \"application/json\" }),\n    });\n    const { data } = await response.json();\n    return data?.clients;\n  } catch (error) {\n    console.error(error);\n    return [];\n  }\n};\nconst fetchMembershipPlans = async () => {\n  try {\n    const response = await fetch(\"/memberships\", {\n      headers: new Headers({ \"content-type\": \"application/json\" }),\n    });\n    const { data } = await response.json();\n    return data?.membershipPlans;\n  } catch (error) {\n    console.error(error);\n    return [];\n  }\n};\nconst subscribe = async (form = {}, onSuccess = () => {}) => {\n  try {\n    const response = await fetch(\"/members\", {\n      method: \"POST\",\n      body: JSON.stringify(form),\n      headers: new Headers({\n        \"content-type\": \"application/json\",\n        \"X-CSRF-Token\": window.csrf,\n      }),\n    });\n    if (response.status === 200) {\n      onSuccess();\n    }\n  } catch (error) {\n    console.error(error);\n  } finally {\n    $(\"#subscribeClientModal\").modal(\"hide\");\n  }\n};\n\nconst cancelSubscription = async (id = 0, onSuccess = () => {}) => {\n  try {\n    const response = await fetch(`/subscriptions/${id}`, {\n      method: \"DELETE\",\n      headers: new Headers({\n        \"content-type\": \"application/json\",\n        \"X-CSRF-Token\": window.csrf,\n      }),\n    });\n    if (response.status === 200) {\n      onSuccess();\n    }\n  } catch (error) {\n    console.error(error);\n  }\n};\n\n\n//# sourceURL=webpack:///./views/admin/members/fetch.js?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The require scope
/******/ 	var __webpack_require__ = {};
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = {};
/******/ 	__webpack_modules__["./views/admin/members/fetch.js"](0, __webpack_exports__, __webpack_require__);
/******/ 	
/******/ })()
;