"use strict"
// Written Docs for this tutorial step can be found here:
// https://github.com/SAP/chevrotain/blob/master/docs/tutorial/step2_parsing.md

// Tutorial Step 2:

// Adding a Parser (grammar only, only reads the input without any actions).
// Using the Token Vocabulary defined in the previous step.

const selectLexer = require("./query-lexer")
const Parser = require("chevrotain").Parser
const tokenVocabulary = selectLexer.tokenVocabulary

// individual imports, prefer ES6 imports if supported in your runtime/transpiler...
const Select = tokenVocabulary.Select
const From = tokenVocabulary.From
const Where = tokenVocabulary.Where
const Identifier = tokenVocabulary.Identifier
const Integer = tokenVocabulary.Integer
const GreaterThan = tokenVocabulary.GreaterThan
const LessThan = tokenVocabulary.LessThan
const Comma = tokenVocabulary.Comma

const Due = tokenVocabulary.Due
const PeriodHistoric = tokenVocabulary.PeriodHistoric
const PeriodFuture = tokenVocabulary.PeriodFuture
const PeriodSpecifier = tokenVocabulary.PeriodSpecifier
const PeriodSpecificYesterday = tokenVocabulary.PeriodSpecificYesterday
const PeriodSpecificTomorrow = tokenVocabulary.PeriodSpecificTomorrow
const PeriodSpecificToday = tokenVocabulary.PeriodSpecificToday

// ----------------- parser -----------------
class SelectParser extends Parser {
    constructor(input, config) {
        super(input, tokenVocabulary, config)
        const $ = this

        this.selectStatement = $.RULE("selectStatement", () => {
//          $.OR([
              //{ALT: () => {$.CONSUME(Due)}},
              //{ALT: () => {$.CONSUME(PeriodFuture)}}
          //])

            $.SUBRULE($.dueClause)
  //          $.SUBRULE($.fromClause)
      //      $.OPTION(() => {
    //            $.SUBRULE($.whereClause)
        //    })
        })

        this.dueClause = $.RULE("dueClause", () => {
            $.CONSUME(Due)
            $.SUBRULE($.periodClause)
            //$.OR([
                //{ALT: () => {$.CONSUME(PeriodHistoric)}},
                //{ALT: () => {$.CONSUME(PeriodFuture)}}
            //])
        })

        this.periodClause = $.RULE("periodClause", () => {
            $.OR([
                {ALT: () => {$.SUBRULE($.periodClauseHistoric)}},
                {ALT: () => {$.SUBRULE($.periodClauseFuture)}},
                {ALT: () => {$.CONSUME(PeriodSpecificYesterday)}},
                {ALT: () => {$.CONSUME(PeriodSpecificTomorrow)}},
                {ALT: () => {$.CONSUME(PeriodSpecificToday)}}
            ])
        })

        this.periodClauseHistoric = $.RULE("periodClauseHistoric", () => {
            $.CONSUME(PeriodHistoric)
            $.CONSUME(Integer)
            $.CONSUME(PeriodSpecifier)
        })

        this.periodClauseFuture = $.RULE("periodClauseFuture", () => {
            $.CONSUME(PeriodFuture)
            $.CONSUME(Integer)
            $.CONSUME(PeriodSpecifier)
        })


        this.selectClause = $.RULE("selectClause", () => {
            //$.CONSUME(Select)
            $.CONSUME(Due)
            $.AT_LEAST_ONE_SEP({
                SEP: Comma,
                DEF: () => {
                    $.CONSUME(Identifier)
                }
            })
        })

        this.fromClause = $.RULE("fromClause", () => {
            $.CONSUME(From)
            $.CONSUME(Identifier)
        })

        this.whereClause = $.RULE("whereClause", () => {
            $.CONSUME(Where)
            $.SUBRULE($.expression)
        })

        this.expression = $.RULE("expression", () => {
            $.SUBRULE($.atomicExpression)
            $.SUBRULE($.relationalOperator)
            $.SUBRULE2($.atomicExpression) // note the '2' suffix to distinguish
            // from the 'SUBRULE(atomicExpression)'
            // 2 lines above.
        })

        this.atomicExpression = $.RULE("atomicExpression", () => {
            // prettier-ignore
            $.OR([
                {ALT: () => {$.CONSUME(Integer)}},
                {ALT: () => {$.CONSUME(Identifier)}}
            ])
        })

        this.relationalOperator = $.RULE("relationalOperator", () => {
            // prettier-ignore
            $.OR([
                {ALT: () => {$.CONSUME(GreaterThan)}},
                {ALT: () => {$.CONSUME(LessThan)}}
            ])
        })

        // very important to call this after all the rules have been defined.
        // otherwise the parser may not work correctly as it will lack information
        // derived during the self analysis phase.
        Parser.performSelfAnalysis(this)
    }
}

// We only ever need one as the parser internal state is reset for each new input.
const parserInstance = new SelectParser([])

module.exports = {
    parserInstance: parserInstance,

    SelectParser: SelectParser,

    parse: function(inputText) {
        let lexResult = selectLexer.lex(inputText)

        // ".input" is a setter which will reset the parser's internal's state.
        parserInstance.input = lexResult.tokens

        // No semantic actions so this won't return anything yet.
        parserInstance.selectStatement()

        if (parserInstance.errors.length > 0) {
            throw Error(
                "Sad sad panda, parsing errors detected!\n" +
                    parserInstance.errors[0].message
            )
        }
    }
}
