from django import forms

class QueryForm(forms.Form):
    start_date = forms.DateTimeField(
        input_formats=['%d/%m/%Y %H:%M']
    )

    end_date =forms.DateTimeField(
        input_formats=['%d/%m/%Y %H:%M']
    )




