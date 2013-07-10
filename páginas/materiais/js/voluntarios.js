$(function() {
	selecaoDeLider();
	selecaoDeZona();
	confirmacao();
});

function selecaoDeLider() {
	configurarBusca();

	$('#selecionar-lider').click(function(e) {
		e.preventDefault();
		$('#selecao-lider').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
			modal: true,
			title: "Zonas",
			buttons: {
				"Não": function() {
					$(this).dialog("close");
				},
				"Sim": function() {
					$("#input-lider").val($('#lider').val());
					$(this).dialog("close");
				}
			}
		});
	})
}

function configurarBusca() {
	var comprimento = 3;
	var buscaNoServidor = "";
	var buscaAnterior = "";
	var resultados = [];

	$("#busca-lider").keyup(function() {
		var novaBusca = $("#busca-lider").val();

		if (novaBusca.length == comprimento && novaBusca != buscaNoServidor) {
			buscaNoServidor = novaBusca;
			obterLideres(buscaNoServidor);

		} else if (novaBusca.length < comprimento) {
			$("#lideres").html("");

		} else if (novaBusca.length >= comprimento) {
			elemento = novaBusca.length > buscaAnterior.length?
				'#lideres' : resultados;

			var html = $(elemento).filter(function() {
				return $(this).text().contains(novaBusca)
			})

			$('#lideres').html(html);
		}

		buscaAnterior = novaBusca;
	});

	var obterLideres = function(busca) {
		$.get("busca/lideres?contem=" + busca, function(lideres) {
			var lista = [];
			for (lider in lideres) {
				var id = lideres[i].id
				var elemento = "<label for='lider" + id + "'>"
				"<input type='radio' class='lideres' name='lideres' id='lider" + id + "'/>"
				lideres[i].nome + "</label><br/>";

				lista.push(elemento);
			}

			$("#lideres").html(lista);
			resultados = lista;
		});
	}
}

function selecaoDeZona() {
	$('#selecionar-zona').click(function(e) {
		e.preventDefault();
		$('#selecao-zona').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
			modal: true,
			title: "Zonas",
			buttons: {
				"Não": function() {
					$(this).dialog("close");
				},
				"Sim": function() {
					$("#input-zona").val($('#zona').val());
					$(this).dialog("close");
				}
			}
		});
	})
}

function confirmacao() {
	var textoOriginal = $("#texto-turnos").text()

	$('#enviar').click(function(e) {
		var turnos = $.map($("#turnos > input:checked + label"), function(x) {
			return $(x).text()
		});

		if (turnos.length == 0) {
			$("#texto-turnos").text("Você não selecionou nenhum turno!");
		} else {
			$("#texto-turnos").text(textoOriginal)
		}
		$("#confirmacao-turnos").html(turnos.map(function(turno) {return "<li>" + turno + "</li>"}));
		$("#confirmacao-zona").html($("#zona > :selected").html());

		e.preventDefault();
		$('#confirmacao').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
			modal: true,
			title: "Confirmação",
			buttons: {
				"Não": function() {
					$(this).dialog("close");
				},
				"Sim": function() {
					$("#formulario").submit();
					$(this).dialog("close");
				}
			}
		});
	});
}
