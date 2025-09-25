import React, { useState } from "react";
import "./App.scss";
import "./styles/Global.css";
import Node from "./components/Canva/Node";
import NavbarComponent from "./components/Navbar/NavbarComponent";
import AboutModal from "./components/About/AboutModal";

const App: React.FC = () => {
  const [showAbout, setShowAbout] = useState<boolean>(false);

  const toggleAboutPopup = () => {
    setShowAbout(!showAbout);
  };

  return (
    <div className="app-container">
      <NavbarComponent version={""} />

      <Node />
      <AboutModal showAbout={showAbout} toggleAboutPopup={toggleAboutPopup} />
    </div>
  );
};

export default App;
