describe("template spec", () => {
  it("should enter email and password if the login form exists and submit", () => {
    cy.visit("http://localhost:3000/");
    cy.get('[data_testid="loginBtn"]').should("exist").click();
    cy.get('[data_testid="toggleLoginRegister"]').should("exist").click();
    cy.get('[data_testid="registerName"]').should("exist").type('Shubham Singh');
    cy.get('[data_testid="registerEmail"]').should("exist").type('shu@ufl.edu');
    cy.get('[data_testid="registerPassword"]').should("exist").type('Shubham16!');
    cy.get('[data_testid="submitRegisterBtn"]').should("exist").click();
  });  
});
