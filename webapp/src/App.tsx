import React, { useState } from "react";
import "./App.scss";
import "./styles/Global.css";
import Node from "./components/Canva/Node";
import NavbarComponent from "./components/Navbar/NavbarComponent";
import AboutModal from "./components/About/AboutModal";

/**
 * App é o componente raiz da aplicação.
 * Ele renderiza a barra de navegação, o componente principal de visualização (Node)
 * e o modal "Sobre".
 */
const App: React.FC = () => {
  // Estado para controlar a visibilidade do modal "Sobre".
  const [showAbout, setShowAbout] = useState<boolean>(false);

  // Função para alternar a visibilidade do modal "Sobre".
  const toggleAboutPopup = () => {
    setShowAbout(!showAbout);
  };

  return (
    <div className="app-container">
      {/* A barra de navegação superior. A propriedade 'version' não está sendo usada no momento. */}
      <NavbarComponent version={""} />

      {/* O componente principal que contém a visualização ReactFlow e os painéis de dados. */}
      <Node />

      {/* O modal "Sobre", que é exibido quando showAbout é verdadeiro. */}
      <AboutModal showAbout={showAbout} toggleAboutPopup={toggleAboutPopup} />
    </div>
  );
};

export default App;
